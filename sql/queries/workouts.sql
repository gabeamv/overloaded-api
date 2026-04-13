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

-- name: GetAllWorkoutsByUserId :many
SELECT * FROM workouts
WHERE user_id = $1;

-- name: GetWorkoutById :one
SELECT * FROM workouts
WHERE id = $1;

-- name: DeleteWorkoutById :exec
DELETE FROM workouts
WHERE id = $1;