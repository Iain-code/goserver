-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: FindUserEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdatePassEmail :exec
UPDATE users
SET email = $2, hashed_password = $3
WHERE users.id = $1;