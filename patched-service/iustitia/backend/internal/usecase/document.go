package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	logkit "github.com/wahrwelt-kit/go-logkit"

	"github.com/TakuyaYagam1/iustitia/internal/apperr"
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/repo"
	sqlc "github.com/TakuyaYagam1/iustitia/internal/repo/sqlc"
	"github.com/TakuyaYagam1/iustitia/pkg/docgen"
)

type Document struct {
	store  repo.Store
	logger logkit.Logger
}

func NewDocument(store repo.Store, logger logkit.Logger) *Document {
	return &Document{store: store, logger: logger}
}

func (u *Document) Generate(ctx context.Context, caseID, authorID uuid.UUID, userTemplate string) (*domain.Document, error) {
	if userTemplate == "" {
		return nil, apperr.ErrBadRequest
	}

	caseRow, err := u.store.GetCaseByID(ctx, caseID.String())
	if err != nil {
		return nil, apperr.ErrNotFound
	}

	verdict := ""
	if caseRow.Verdict != nil {
		verdict = *caseRow.Verdict
	}
	if verdict == "" {
		// Для ещё не вынесенных приговоров подставляем placeholder чтобы
		// шаблон verdict.tpl не оставлял пустую строку после постановки.
		verdict = "Вопрос о мере наказания будет разрешён при постановке приговора."
	}
	docCtx := &docgen.DocumentContext{
		CaseID:    caseRow.ID,
		Defendant: caseRow.Defendant,
		Verdict:   verdict,
		Details:   caseRow.Crime,
	}

	// Резолвим template_name через canonical registry: если это один из
	// известных пресетов (summons/indictment/verdict) - берём готовое тело
	// из registry. Иначе (для чекера и произвольных тел) - используем
	// userTemplate как есть. Либо случай всё равно проходит через
	// docgen.Generate, который безопасен (whitelist-плейсхолдеры, без
	// text/template - см. F3).
	tplBody, isCanonical := docgen.ResolveTemplate(userTemplate)
	if !isCanonical {
		tplBody = userTemplate
	}

	content, err := docgen.Generate(tplBody, docCtx)
	if err != nil {
		return nil, fmt.Errorf("Document - Generate - docgen.Generate: %w", err)
	}

	id := uuid.New()
	row, err := u.store.CreateDocument(ctx, sqlc.CreateDocumentParams{
		ID:       id.String(),
		CaseID:   caseID.String(),
		AuthorID: authorID.String(),
		Content:  content,
		Template: userTemplate,
	})
	if err != nil {
		return nil, fmt.Errorf("Document - Generate - store.CreateDocument: %w", err)
	}
	return documentFromRow(row)
}

func (u *Document) ListByCase(ctx context.Context, caseID uuid.UUID) ([]*domain.Document, error) {
	rows, err := u.store.ListDocumentsByCase(ctx, caseID.String())
	if err != nil {
		return nil, fmt.Errorf("Document - ListByCase - store.ListDocumentsByCase: %w", err)
	}
	out := make([]*domain.Document, 0, len(rows))
	for _, r := range rows {
		d, err := documentFromRow(r)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

func (u *Document) GetByID(ctx context.Context, id uuid.UUID) (*domain.Document, error) {
	row, err := u.store.GetDocumentByID(ctx, id.String())
	if err != nil {
		return nil, apperr.ErrNotFound
	}
	return documentFromRow(row)
}

func documentFromRow(r sqlc.Document) (*domain.Document, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, errors.New("usecase - documentFromRow - malformed id")
	}
	caseID, err := uuid.Parse(r.CaseID)
	if err != nil {
		return nil, errors.New("usecase - documentFromRow - malformed case_id")
	}
	authorID, err := uuid.Parse(r.AuthorID)
	if err != nil {
		return nil, errors.New("usecase - documentFromRow - malformed author_id")
	}
	return &domain.Document{
		ID:        id,
		CaseID:    caseID,
		AuthorID:  authorID,
		Content:   r.Content,
		Template:  r.Template,
		CreatedAt: r.CreatedAt,
	}, nil
}
