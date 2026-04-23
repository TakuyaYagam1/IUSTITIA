package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	logkit "github.com/wahrwelt-kit/go-logkit"

	"github.com/TakuyaYagam1/iustitia/internal/apperr"
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/repo"
	sqlc "github.com/TakuyaYagam1/iustitia/internal/repo/sqlc"
)

type Case struct {
	store  repo.Store
	logger logkit.Logger
}

func NewCase(store repo.Store, logger logkit.Logger) *Case {
	return &Case{store: store, logger: logger}
}

func (u *Case) List(ctx context.Context, limit, offset int) ([]*domain.Case, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := u.store.ListCases(ctx, sqlc.ListCasesParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("Case - List - store.ListCases: %w", err)
	}
	out := make([]*domain.Case, 0, len(rows))
	for _, row := range rows {
		c, err := caseFromRow(row)
		if err != nil {
			return nil, fmt.Errorf("Case - List - caseFromRow: %w", err)
		}
		out = append(out, c)
	}
	return out, nil
}

func (u *Case) GetByID(ctx context.Context, id uuid.UUID) (*domain.Case, error) {
	row, err := u.store.GetCaseByID(ctx, id.String())
	if err != nil {
		return nil, apperr.ErrNotFound
	}
	return caseFromRow(row)
}

type SearchRequest struct {
	Q         string
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

func (u *Case) Search(ctx context.Context, req SearchRequest) ([]*domain.Case, error) {
	orderBy := req.OrderBy
	direction := req.Direction

	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 200 {
		req.Limit = 200
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	rows, err := u.store.SearchCases(ctx, repo.SearchCasesRequest{
		Q:         req.Q,
		OrderBy:   orderBy,
		Direction: direction,
		Limit:     req.Limit,
		Offset:    req.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("Case - Search - store.SearchCases: %w", err)
	}
	out := make([]*domain.Case, 0, len(rows))
	for _, row := range rows {
		c, err := caseFromRow(row)
		if err != nil {
			return nil, fmt.Errorf("Case - Search - caseFromRow: %w", err)
		}
		out = append(out, c)
	}
	return out, nil
}

type VerdictResult struct {
	Case         *domain.Case
	Document     *domain.Document
	ArchiveEntry *domain.ArchiveEntry
}

type HearingItem struct {
	Case    *domain.Case
	Opinion *domain.CaseOpinion
}

func (u *Case) Create(ctx context.Context, citizenID uuid.UUID, defendant, crime, firstText string) (*domain.Case, error) {
	defendant = strings.TrimSpace(defendant)
	crime = strings.TrimSpace(crime)
	firstText = strings.TrimSpace(firstText)
	if defendant == "" || crime == "" || firstText == "" {
		return nil, apperr.ErrBadRequest
	}

	caseID := uuid.New()
	complaintID := uuid.New()
	authorStr := citizenID.String()

	var created sqlc.Case
	err := u.store.WithTx(ctx, func(s repo.Store) error {
		row, err := s.CreateCase(ctx, sqlc.CreateCaseParams{
			ID:        caseID.String(),
			Defendant: defendant,
			Crime:     crime,
			AuthorID:  &authorStr,
		})
		if err != nil {
			return fmt.Errorf("CreateCase: %w", err)
		}
		if _, err := s.CreateComplaint(ctx, sqlc.CreateComplaintParams{
			ID:       complaintID.String(),
			CaseID:   caseID.String(),
			AuthorID: citizenID.String(),
			Text:     firstText,
		}); err != nil {
			return fmt.Errorf("CreateComplaint: %w", err)
		}
		created = row
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Case - Create - tx: %w", err)
	}
	return caseFromRow(created)
}

func (u *Case) Accept(ctx context.Context, caseID, prosecutorID uuid.UUID) (*domain.Case, error) {
	userRow, err := u.store.GetUserByID(ctx, prosecutorID.String())
	if err != nil {
		return nil, apperr.ErrBadRequest
	}
	if domain.Role(userRow.Role) != domain.RoleProsecutor {
		return nil, apperr.ErrBadRequest
	}
	prosStr := prosecutorID.String()
	row, err := u.store.AcceptCase(ctx, sqlc.AcceptCaseParams{
		AssignedProsecutorID: &prosStr,
		ID:                   caseID.String(),
	})
	if err != nil {
		return nil, apperr.ErrConflict
	}
	return caseFromRow(row)
}

func (u *Case) Dismiss(ctx context.Context, caseID uuid.UUID, reason string) (*VerdictResult, error) {
	_ = strings.TrimSpace(reason)
	var res VerdictResult
	err := u.store.WithTx(ctx, func(s repo.Store) error {
		row, err := s.DismissCase(ctx, caseID.String())
		if err != nil {
			return apperr.ErrConflict
		}
		archID := uuid.New()
		caseIDStr := row.ID
		archRow, err := s.CreateArchiveEntry(ctx, sqlc.CreateArchiveEntryParams{
			ID:             archID.String(),
			CaseID:         &caseIDStr,
			Defendant:      row.Defendant,
			FinalVerdict:   string(domain.VerdictDismissed),
			Sentence:       nil,
			ClassifiedNote: row.ClassifiedNote,
		})
		if err != nil {
			return fmt.Errorf("CreateArchiveEntry: %w", err)
		}
		c, err := caseFromRow(row)
		if err != nil {
			return err
		}
		ae, err := archiveFromRow(archRow)
		if err != nil {
			return err
		}
		res.Case = c
		res.ArchiveEntry = ae
		return nil
	})
	if err != nil {
		if errors.Is(err, apperr.ErrConflict) {
			return nil, err
		}
		return nil, fmt.Errorf("Case - Dismiss - tx: %w", err)
	}
	return &res, nil
}

func (u *Case) FileOpinion(
	ctx context.Context,
	caseID, prosecutorID uuid.UUID,
	prelim domain.Verdict,
	reasoning string,
) (*domain.CaseOpinion, error) {
	if !prelim.Valid() {
		return nil, apperr.ErrBadRequest
	}
	reasoning = strings.TrimSpace(reasoning)
	if reasoning == "" {
		return nil, apperr.ErrBadRequest
	}

	opinionID := uuid.New()
	prosStr := prosecutorID.String()
	var created sqlc.CaseOpinion
	err := u.store.WithTx(ctx, func(s repo.Store) error {
		row, err := s.CreateCaseOpinion(ctx, sqlc.CreateCaseOpinionParams{
			ID:                 opinionID.String(),
			CaseID:             caseID.String(),
			ProsecutorID:       prosecutorID.String(),
			PreliminaryVerdict: string(prelim),
			Reasoning:          reasoning,
		})
		if err != nil {
			return apperr.ErrConflict
		}
		if _, err := s.MarkCaseHearing(ctx, sqlc.MarkCaseHearingParams{
			ID:                   caseID.String(),
			AssignedProsecutorID: &prosStr,
		}); err != nil {
			return apperr.ErrForbidden
		}
		created = row
		return nil
	})
	if err != nil {
		if errors.Is(err, apperr.ErrConflict) || errors.Is(err, apperr.ErrForbidden) {
			return nil, err
		}
		return nil, fmt.Errorf("Case - FileOpinion - tx: %w", err)
	}
	return opinionFromRow(created)
}

func (u *Case) GetOpinion(ctx context.Context, caseID uuid.UUID) (*domain.CaseOpinion, error) {
	row, err := u.store.GetCaseOpinionByCase(ctx, caseID.String())
	if err != nil {
		return nil, apperr.ErrNotFound
	}
	return opinionFromRow(row)
}

func (u *Case) ListHearings(ctx context.Context) ([]*HearingItem, error) {
	rows, err := u.store.ListCasesByStatus(ctx, string(domain.CaseStatusHearing))
	if err != nil {
		return nil, fmt.Errorf("Case - ListHearings - store.ListCasesByStatus: %w", err)
	}
	out := make([]*HearingItem, 0, len(rows))
	for _, r := range rows {
		c, err := caseFromRow(r)
		if err != nil {
			return nil, fmt.Errorf("Case - ListHearings - caseFromRow: %w", err)
		}
		opRow, err := u.store.GetCaseOpinionByCase(ctx, r.ID)
		if err != nil {
			out = append(out, &HearingItem{Case: c, Opinion: nil})
			continue
		}
		op, err := opinionFromRow(opRow)
		if err != nil {
			return nil, fmt.Errorf("Case - ListHearings - opinionFromRow: %w", err)
		}
		out = append(out, &HearingItem{Case: c, Opinion: op})
	}
	return out, nil
}

type docgenRenderer interface {
	RenderVerdict(caseSeq int64, defendant, crime, details, verdict string) (string, error)
}

func (u *Case) FinalizeVerdict(
	ctx context.Context,
	caseID, judgeID uuid.UUID,
	verdict domain.Verdict,
	sentence *string,
	reasoning string,
	renderer docgenRenderer,
) (*VerdictResult, error) {
	if !verdict.Valid() {
		return nil, apperr.ErrBadRequest
	}
	reasoning = strings.TrimSpace(reasoning)
	if reasoning == "" {
		return nil, apperr.ErrBadRequest
	}
	if verdict == domain.VerdictGuilty {
		if sentence == nil || strings.TrimSpace(*sentence) == "" {
			return nil, apperr.ErrBadRequest
		}
	} else {
		sentence = nil
	}

	caseRow, err := u.store.GetCaseByID(ctx, caseID.String())
	if err != nil {
		return nil, apperr.ErrNotFound
	}
	if domain.CaseStatus(caseRow.Status) != domain.CaseStatusHearing {
		return nil, apperr.ErrConflict
	}

	opRow, err := u.store.GetCaseOpinionByCase(ctx, caseID.String())
	if err != nil {
		return nil, apperr.ErrConflict
	}

	details := fmt.Sprintf(
		"%s\n\nЗаключение прокурора: %s.\nОбоснование прокурора: %s\n\nОбоснование суда: %s",
		caseRow.Crime, opRow.PreliminaryVerdict, opRow.Reasoning, reasoning,
	)
	verdictRu := verdictToRu(verdict)
	if sentence != nil {
		verdictRu = verdictRu + ". Мера наказания: " + *sentence
	}
	content, err := renderer.RenderVerdict(caseRow.SeqNum, caseRow.Defendant, caseRow.Crime, details, verdictRu)
	if err != nil {
		return nil, fmt.Errorf("Case - FinalizeVerdict - renderer: %w", err)
	}

	var res VerdictResult
	err = u.store.WithTx(ctx, func(s repo.Store) error {
		docID := uuid.New()
		docRow, err := s.CreateDocument(ctx, sqlc.CreateDocumentParams{
			ID:       docID.String(),
			CaseID:   caseID.String(),
			AuthorID: judgeID.String(),
			Content:  content,
			Template: "verdict",
		})
		if err != nil {
			return fmt.Errorf("CreateDocument: %w", err)
		}
		vStr := string(verdict)
		updated, err := s.FinalizeCaseVerdict(ctx, sqlc.FinalizeCaseVerdictParams{
			Verdict: &vStr,
			ID:      caseID.String(),
		})
		if err != nil {
			return apperr.ErrConflict
		}
		archID := uuid.New()
		caseIDStr := caseID.String()
		archRow, err := s.CreateArchiveEntry(ctx, sqlc.CreateArchiveEntryParams{
			ID:             archID.String(),
			CaseID:         &caseIDStr,
			Defendant:      updated.Defendant,
			FinalVerdict:   vStr,
			Sentence:       sentence,
			ClassifiedNote: updated.ClassifiedNote,
		})
		if err != nil {
			return fmt.Errorf("CreateArchiveEntry: %w", err)
		}

		c, err := caseFromRow(updated)
		if err != nil {
			return err
		}
		doc, err := documentFromRow(docRow)
		if err != nil {
			return err
		}
		ae, err := archiveFromRow(archRow)
		if err != nil {
			return err
		}
		res.Case = c
		res.Document = doc
		res.ArchiveEntry = ae
		return nil
	})
	if err != nil {
		if errors.Is(err, apperr.ErrConflict) {
			return nil, err
		}
		return nil, fmt.Errorf("Case - FinalizeVerdict - tx: %w", err)
	}
	return &res, nil
}

func verdictToRu(v domain.Verdict) string {
	switch v {
	case domain.VerdictGuilty:
		return "Признать подсудимого виновным"
	case domain.VerdictAcquitted:
		return "Признать подсудимого невиновным"
	case domain.VerdictDismissed:
		return "Производство по делу прекратить"
	}
	return string(v)
}

func caseFromRow(r sqlc.Case) (*domain.Case, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, errors.New("usecase - caseFromRow - malformed id")
	}
	c := &domain.Case{
		ID:             id,
		SeqNum:         r.SeqNum,
		Defendant:      r.Defendant,
		Crime:          r.Crime,
		Status:         domain.CaseStatus(r.Status),
		ClassifiedNote: r.ClassifiedNote,
		CreatedAt:      r.CreatedAt,
	}
	if r.Verdict != nil {
		v := domain.Verdict(*r.Verdict)
		c.Verdict = &v
	}
	if r.AuthorID != nil {
		if aid, err := uuid.Parse(*r.AuthorID); err == nil {
			c.AuthorID = &aid
		}
	}
	if r.AssignedProsecutorID != nil {
		if pid, err := uuid.Parse(*r.AssignedProsecutorID); err == nil {
			c.AssignedProsecutorID = &pid
		}
	}
	return c, nil
}

func opinionFromRow(r sqlc.CaseOpinion) (*domain.CaseOpinion, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, errors.New("usecase - opinionFromRow - malformed id")
	}
	caseID, err := uuid.Parse(r.CaseID)
	if err != nil {
		return nil, errors.New("usecase - opinionFromRow - malformed case_id")
	}
	prosID, err := uuid.Parse(r.ProsecutorID)
	if err != nil {
		return nil, errors.New("usecase - opinionFromRow - malformed prosecutor_id")
	}
	return &domain.CaseOpinion{
		ID:                 id,
		CaseID:             caseID,
		ProsecutorID:       prosID,
		PreliminaryVerdict: domain.Verdict(r.PreliminaryVerdict),
		Reasoning:          r.Reasoning,
		FiledAt:            r.FiledAt,
	}, nil
}
