-- name: CreateExercise :one
INSERT INTO exercises (name)
VALUES (
    $1
)
RETURNING *;