-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1 FOR NO KEY UPDATE;

-- name: CreateUser :one
INSERT INTO users (username, hashed_password, full_name, email)
VALUES ($1, $2, $3, $4)
RETURNING *;