-- name: GetDirectiveByCode :one
SELECT id, directive_code, secret_payload, classification, issued_at
FROM mtb_directives
WHERE directive_code = ?
LIMIT 1;

-- name: ListDirectivesPublic :many
SELECT id, directive_code, secret_payload, classification, issued_at
FROM mtb_directives
WHERE classification = 'public'
ORDER BY issued_at DESC
LIMIT ? OFFSET ?;
