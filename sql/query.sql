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

-- name: GetUserBySessionID :one
SELECT * FROM users
WHERE session_id = $1 AND deleted_at IS NULL AND updated_at + interval '10' hour > now()
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  email, verification_code, verification_at
) VALUES (
  $1, $2, now()
)
RETURNING *;

-- name: UpdateVerificationCode :exec
UPDATE users 
SET verification_code = $2, verification_at = now()
WHERE id = $1;

-- name: UpdateSessionID :exec
UPDATE users 
SET session_id = $2, updated_at = now()
WHERE id = $1; 
