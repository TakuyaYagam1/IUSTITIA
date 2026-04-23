package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	logkit "github.com/wahrwelt-kit/go-logkit"

	"github.com/TakuyaYagam1/iustitia/internal/apperr"
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/repo"
	sqlc "github.com/TakuyaYagam1/iustitia/internal/repo/sqlc"
)

const maxEvidenceBytes = 5 << 20

type Complaint struct {
	store      repo.Store
	logger     logkit.Logger
	evidenceHC *http.Client
}

func NewComplaint(store repo.Store, logger logkit.Logger) *Complaint {
	tr := &http.Transport{}
	tr.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	return &Complaint{
		store:  store,
		logger: logger,
		evidenceHC: &http.Client{
			Transport: tr,
			Timeout:   10 * time.Second,
		},
	}
}

func (u *Complaint) Create(ctx context.Context, caseID, authorID uuid.UUID, text string) (*domain.Complaint, error) {
	if text == "" {
		return nil, apperr.ErrBadRequest
	}
	id := uuid.New()
	row, err := u.store.CreateComplaint(ctx, sqlc.CreateComplaintParams{
		ID:       id.String(),
		CaseID:   caseID.String(),
		AuthorID: authorID.String(),
		Text:     text,
	})
	if err != nil {
		return nil, fmt.Errorf("Complaint - Create - store.CreateComplaint: %w", err)
	}
	return complaintFromRow(row)
}

func (u *Complaint) AttachEvidence(ctx context.Context, complaintID uuid.UUID, url string) (*domain.Complaint, error) {
	if url == "" {
		return nil, apperr.ErrBadRequest
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	urlCopy := url
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
