-- name: GetUserByUsername :one
SELECT id, username, password, role, dome, created_at
FROM users
WHERE username = ?
LIMIT 1;

-- name: GetUserByID :one
SELECT id, username, password, role, dome, created_at
FROM users
WHERE id = ?
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, username, password, role, dome)
VALUES (?, ?, ?, ?, ?)
RETURNING id, username, password, role, dome, created_at;

-- name: ListUsersByRole :many
SELECT id, username, password, role, dome, created_at
FROM users
WHERE role = ?
ORDER BY username;
