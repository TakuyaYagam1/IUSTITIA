-- name: CreateDocument :one
INSERT INTO documents (id, case_id, author_id, content, template)
VALUES (?, ?, ?, ?, ?)
RETURNING id, case_id, author_id, content, template, created_at;

-- name: GetDocumentByID :one
SELECT id, case_id, author_id, content, template, created_at
FROM documents
WHERE id = ?
LIMIT 1;

-- name: ListDocumentsByCase :many
SELECT id, case_id, author_id, content, template, created_at
FROM documents
WHERE case_id = ?
ORDER BY created_at DESC;
