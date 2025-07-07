-- name: CreateURL :one
INSERT INTO urls (code, original_url) VALUES ($1, $2)
RETURNING id, code, original_url, created_at;

-- name: GetURLByCode :one
SELECT id, code, original_url, created_at FROM urls WHERE code = $1;