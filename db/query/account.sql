-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1
LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY owner
LIMIT $1
OFFSET $2
FOR NO KEY UPDATE;

-- name: CreateAccount :one
INSERT INTO accounts (owner, balance, currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateAccount :exec
UPDATE accounts
set owner = $2, balance = $3, currency = $4
WHERE id = $1;

-- name: UpdateAccountBalance :one
UPDATE accounts
set balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;