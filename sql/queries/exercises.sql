-- name: CreateExercise :one
INSERT INTO exercises (name, is_custom, user_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;