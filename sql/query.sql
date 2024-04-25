-- name: GetAPIKey :one
SELECT * FROM api_keys WHERE api_key = ? AND deleted_at IS NULL LIMIT 1;
