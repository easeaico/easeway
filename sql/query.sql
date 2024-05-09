-- name: GetAPIKey :one
SELECT * FROM api_keys 
WHERE key = $1 AND deleted_at IS NULL 
LIMIT 1;

-- name: CreateAPIKey :one
INSERT INTO api_keys (
  user_id, name, key
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: CreateOutcome :one
INSERT INTO outcomes (
  user_id, key_id, prompt_tokens, completion_tokens, total_tokens, rt, fee_rate, cost
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;


-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  email, verification_code
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateVerificationCode :exec
UPDATE users 
SET verification_code = $2
WHERE id = $1;

-- name: UpdateSessionID :exec
UPDATE users 
SET session_id = $2
WHERE id = $1; 
