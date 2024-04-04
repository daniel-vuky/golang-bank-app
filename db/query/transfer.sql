-- name: GetTransfer :one
SELECT *
FROM transfers
WHERE id = $1
LIMIT 1 FOR NO KEY UPDATE;

-- name: ListTransfers :many
SELECT *
FROM transfers
WHERE from_account_id = $1
   or to_account_id = $2
ORDER BY id DESC
LIMIT $3 OFFSET $4 FOR NO KEY UPDATE;

-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateTransfer :exec
UPDATE transfers
set amount = $2
WHERE id = $1;

-- name: DeleteTransfer :exec
DELETE
FROM transfers
WHERE id = $1;