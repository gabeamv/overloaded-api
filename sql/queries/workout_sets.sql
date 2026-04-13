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