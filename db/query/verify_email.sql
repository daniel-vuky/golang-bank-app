-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
    username,
    email,
    token
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET is_used = true
WHERE id = $1 AND
      token = $2 AND
      is_used = false AND
      expired_at > NOW()
RETURNING *;