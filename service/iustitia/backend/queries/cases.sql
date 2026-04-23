-- name: ListCases :many
SELECT id, seq_num, defendant, crime, status, verdict, classified_note, author_id, assigned_prosecutor_id, created_at
FROM cases
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetCaseByID :one
SELECT id, seq_num, defendant, crime, status, verdict, classified_note, author_id, assigned_prosecutor_id, created_at
FROM cases
WHERE id = ?
LIMIT 1;


-- name: CreateCase :one
INSERT INTO cases (id, seq_num, defendant, crime, status, author_id)
VALUES (?, (SELECT COALESCE(MAX(seq_num), 0) + 1 FROM cases), ?, ?, 'draft', ?)
RETURNING id, seq_num, defendant, crime, status, verdict, classified_note, author_id, assigned_prosecutor_id, created_at;

-- name: AcceptCase :one
UPDATE cases
SET status                 = 'assigned',
    assigned_prosecutor_id = ?
WHERE id = ? AND status IN ('draft', 'open')
RETURNING id, seq_num, defendant, crime, status, verdict, classified_note, author_id, assigned_prosecutor_id, created_at;

-- name: DismissCase :one
UPDATE cases
SET status  = 'closed',
    verdict = 'dismissed'
WHERE id = ? AND status IN ('draft', 'open')
RETURNING id, seq_num, defendant, crime, status, verdict, classified_note, author_id, assigned_prosecutor_id, created_at;

-- name: MarkCaseHearing :one
UPDATE cases
SET status = 'hearing'
WHERE id = ? AND status = 'assigned' AND assigned_prosecutor_id = ?
RETURNING id, seq_num, defendant, crime, status, verdict, classified_note, author_id, assigned_prosecutor_id, created_at;

-- name: FinalizeCaseVerdict :one
UPDATE cases
SET status  = 'closed',
    verdict = ?
WHERE id = ? AND status = 'hearing'
RETURNING id, seq_num, defendant, crime, status, verdict, classified_note, author_id, assigned_prosecutor_id, created_at;

-- name: ListCasesByStatus :many
SELECT id, seq_num, defendant, crime, status, verdict, classified_note, author_id, assigned_prosecutor_id, created_at
FROM cases
WHERE status = ?
ORDER BY created_at DESC;

-- name: ListCasesByAssignee :many
SELECT id, seq_num, defendant, crime, status, verdict, classified_note, author_id, assigned_prosecutor_id, created_at
FROM cases
WHERE assigned_prosecutor_id = ?
  AND status IN ('assigned', 'hearing')
ORDER BY created_at DESC;
