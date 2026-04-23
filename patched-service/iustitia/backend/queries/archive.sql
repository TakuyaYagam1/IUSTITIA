-- name: ListArchive :many
SELECT id, case_id, defendant, final_verdict, sentence, classified_note, archived_at
FROM archive
ORDER BY archived_at DESC
LIMIT ? OFFSET ?;

-- name: GetArchiveByID :one
SELECT id, case_id, defendant, final_verdict, sentence, classified_note, archived_at
FROM archive
WHERE id = ?
LIMIT 1;

-- name: UpdateArchiveSafe :one
UPDATE archive
SET sentence      = COALESCE(?, sentence),
    final_verdict = COALESCE(?, final_verdict)
WHERE id = ?
RETURNING id, case_id, defendant, final_verdict, sentence, classified_note, archived_at;

-- name: CreateArchiveEntry :one
INSERT INTO archive (id, case_id, defendant, final_verdict, sentence, classified_note)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id, case_id, defendant, final_verdict, sentence, classified_note, archived_at;
