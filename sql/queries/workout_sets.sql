-- name: CreateWorkoutSetByWorkoutId :one
INSERT INTO workout_sets (workout_id, exercise_id, progress_track, weight_in_lbs, reps, time_in_seconds)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetSetsByWorkoutId :many
SELECT * FROM workout_sets
WHERE workout_id = $1;

-- name: GetSetById :one
SELECT * FROM workout_sets
WHERE id = $1;

-- name: UpdateSetById :one
UPDATE workout_sets
SET exercise_id = $2, progress_track = $3, weight_in_lbs = $4, reps = $5, time_in_seconds = $6
WHERE id = $1
RETURNING *;

-- name: DeleteSetById :exec
DELETE FROM workout_sets
WHERE id = $1;