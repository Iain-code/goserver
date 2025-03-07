-- name: MakeRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: GetUserFromToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1 
  AND refresh_tokens.revoked_at IS NULL 
  AND refresh_tokens.expires_at > NOW();

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE refresh_tokens.token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = Now(), updated_at = Now()
WHERE refresh_tokens.token = $1;