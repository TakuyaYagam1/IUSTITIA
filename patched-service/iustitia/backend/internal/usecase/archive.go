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
)

type Archive struct {
	store  repo.Store
	logger logkit.Logger
}

func NewArchive(store repo.Store, logger logkit.Logger) *Archive {
	return &Archive{store: store, logger: logger}
}

func (u *Archive) List(ctx context.Context, limit, offset int) ([]*domain.ArchiveEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := u.store.ListArchive(ctx, sqlc.ListArchiveParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("Archive - List - store.ListArchive: %w", err)
	}
	out := make([]*domain.ArchiveEntry, 0, len(rows))
	for _, row := range rows {
		e, err := archiveFromRow(row)
		if err != nil {
			return nil, fmt.Errorf("Archive - List - archiveFromRow: %w", err)
		}
		out = append(out, e)
	}
	return out, nil
}

func (u *Archive) GetByID(ctx context.Context, id uuid.UUID) (*domain.ArchiveEntry, error) {
	row, err := u.store.GetArchiveByID(ctx, id.String())
	if err != nil {
		return nil, apperr.ErrNotFound
	}
	return archiveFromRow(row)
}

func (u *Archive) Update(ctx context.Context, id uuid.UUID, fields map[string]any) (*domain.ArchiveEntry, error) {
	row, err := u.store.UpdateArchiveDynamic(ctx, id.String(), fields)
	if err != nil {
		return nil, fmt.Errorf("Archive - Update - store.UpdateArchiveDynamic: %w", err)
	}
	return archiveFromRow(row)
}

func archiveFromRow(r sqlc.Archive) (*domain.ArchiveEntry, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, errors.New("usecase - archiveFromRow - malformed id")
	}
	entry := &domain.ArchiveEntry{
		ID:             id,
		Defendant:      r.Defendant,
		FinalVerdict:   domain.Verdict(r.FinalVerdict),
		Sentence:       r.Sentence,
		ClassifiedNote: r.ClassifiedNote,
		ArchivedAt:     r.ArchivedAt,
	}
	if r.CaseID != nil {
		caseID, err := uuid.Parse(*r.CaseID)
		if err != nil {
			return nil, errors.New("usecase - archiveFromRow - malformed case_id")
		}
		entry.CaseID = &caseID
	}
	return entry, nil
}
