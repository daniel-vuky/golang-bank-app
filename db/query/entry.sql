-- name: GetEntry :one
SELECT *
FROM entries
WHERE id = $1
LIMIT 1 FOR NO KEY UPDATE;

-- name: ListEntries :many
SELECT *
FROM entries
WHERE account_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3 FOR NO KEY UPDATE;

-- name: CreateEntry :one
INSERT INTO entries (account_id, amount)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateEntry :exec
UPDATE entries
set amount = $2
WHERE id = $1;

-- name: DeleteEntry :exec
DELETE
FROM entries
WHERE id = $1;