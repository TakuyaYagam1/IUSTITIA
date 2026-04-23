-- name: CreateCaseOpinion :one
INSERT INTO case_opinions (id, case_id, prosecutor_id, preliminary_verdict, reasoning)
VALUES (?, ?, ?, ?, ?)
RETURNING id, case_id, prosecutor_id, preliminary_verdict, reasoning, filed_at;

-- name: GetCaseOpinionByCase :one
SELECT id, case_id, prosecutor_id, preliminary_verdict, reasoning, filed_at
FROM case_opinions
WHERE case_id = ?
LIMIT 1;
