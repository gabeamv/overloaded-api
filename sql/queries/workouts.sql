-- name: CreateWorkout :one
INSERT INTO workouts (user_id, started_at, ended_at, volume, prs)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;