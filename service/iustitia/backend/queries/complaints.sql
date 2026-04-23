-- name: CreateComplaint :one
INSERT INTO complaints (id, case_id, author_id, text)
VALUES (?, ?, ?, ?)
RETURNING id, case_id, author_id, text, evidence_url, evidence_data, created_at;

-- name: ListComplaintsByCase :many
SELECT id, case_id, author_id, text, evidence_url, evidence_data, created_at
FROM complaints
WHERE case_id = ?
ORDER BY created_at DESC;

-- name: GetComplaintByID :one
SELECT id, case_id, author_id, text, evidence_url, evidence_data, created_at
FROM complaints
WHERE id = ?
LIMIT 1;

-- name: UpdateComplaintEvidence :one
UPDATE complaints
SET evidence_url  = ?,
    evidence_data = ?
WHERE id = ?
RETURNING id, case_id, author_id, text, evidence_url, evidence_data, created_at;
