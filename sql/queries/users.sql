-- name: CreateUser :one
INSERT INTO users (username, email, hashed_password, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: UpdateUserPasswordById :one
UPDATE users
SET hashed_password = $2, updated_at = $3
WHERE id = $1
RETURNING *;

-- name: UpdateUserUsernameById :one
UPDATE users
SET username = $2, updated_at = $3
WHERE id = $1
RETURNING *;

-- name: UpdateUserEmailById :one
UPDATE users
SET email = $2, updated_at = $3
WHERE id = $1
RETURNING *;