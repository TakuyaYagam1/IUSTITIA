package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	sqlc "github.com/TakuyaYagam1/iustitia/internal/repo/sqlc"
)

type SQLStore struct {
	db *sql.DB
	*sqlc.Queries
}

var _ Store = (*SQLStore)(nil)

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: sqlc.New(db),
	}
}

func (s *SQLStore) WithTx(ctx context.Context, fn func(Store) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("SQLStore - WithTx - BeginTx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	txStore := &SQLStore{
		db:      s.db,
		Queries: s.Queries.WithTx(tx),
	}

	if err := fn(txStore); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("SQLStore - WithTx - Commit: %w", err)
	}
	return nil
}

func (s *SQLStore) SearchCases(ctx context.Context, req SearchCasesRequest) ([]sqlc.Case, error) {
	if req.Limit <= 0 {
		req.Limit = 50
	}

	orderClause := req.OrderBy
	if req.Direction != "" {
		orderClause = orderClause + " " + req.Direction
	}

	qb := sq.Select(
		"id", "seq_num", "defendant", "crime", "status",
		"verdict", "classified_note", "created_at",
	).
		From("cases").
		Where(sq.Like{"defendant": "%" + req.Q + "%"}).
		OrderBy(orderClause).
		Limit(uint64(req.Limit)).
		Offset(uint64(req.Offset))

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("SQLStore - SearchCases - ToSql: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("SQLStore - SearchCases - QueryContext: %w", err)
	}
	defer func() { _ = rows.Close() }()

	cases := make([]sqlc.Case, 0)
	for rows.Next() {
		var c sqlc.Case
		if err := rows.Scan(
			&c.ID, &c.SeqNum, &c.Defendant, &c.Crime, &c.Status,
			&c.Verdict, &c.ClassifiedNote, &c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("SQLStore - SearchCases - Scan: %w", err)
		}
		cases = append(cases, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("SQLStore - SearchCases - rows.Err: %w", err)
	}
	return cases, nil
}

// PATCH - Vuln 4 (Mass Assignment на classified_note).
// Было: classified_note присутствовал в whitelist'е, и Update
// ArchiveDynamic пропускал его в SET-часть SQL - handler передавал его
// из map[string]any напрямую.
// Стало: classified_note убран из whitelist'а. Даже если handler
// (раньше) или чекер теперь пропустит это поле в fields - store его
// молча отбросит.
var archiveMutableColumns = map[string]struct{}{
	"defendant":     {},
	"final_verdict": {},
	"sentence":      {},
}

func (s *SQLStore) UpdateArchiveDynamic(ctx context.Context, id string, fields map[string]any) (sqlc.Archive, error) {
	if id == "" {
		return sqlc.Archive{}, errors.New("SQLStore - UpdateArchiveDynamic - empty id")
	}

	qb := sq.Update("archive")
	applied := 0
	for k, v := range fields {
		if _, ok := archiveMutableColumns[k]; !ok {
			continue
		}
		qb = qb.Set(k, v)
		applied++
	}

	if applied == 0 {
		return s.GetArchiveByID(ctx, id)
	}

	qb = qb.Where(sq.Eq{"id": id})
	query, args, err := qb.ToSql()
	if err != nil {
		return sqlc.Archive{}, fmt.Errorf("SQLStore - UpdateArchiveDynamic - ToSql: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return sqlc.Archive{}, fmt.Errorf("SQLStore - UpdateArchiveDynamic - ExecContext: %w", err)
	}
	return s.GetArchiveByID(ctx, id)
}
