package repo

import (
	"context"

	sqlc "github.com/TakuyaYagam1/iustitia/internal/repo/sqlc"
)

type SearchCasesRequest struct {
	Q         string
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type Store interface {
	sqlc.Querier

	SearchCases(ctx context.Context, req SearchCasesRequest) ([]sqlc.Case, error)
	UpdateArchiveDynamic(ctx context.Context, id string, fields map[string]any) (sqlc.Archive, error)

	WithTx(ctx context.Context, fn func(Store) error) error
}
