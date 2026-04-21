-- name: GetPrOneRepMaxByExerciseIdUserId :one
SELECT COALESCE(
    ROUND(MAX(workout_sets.weight_in_lbs * (1 + (workout_sets.reps / 30.0))), 2),
    0.0
)::DOUBLE PRECISION as orm
FROM workout_sets 
    INNER JOIN workouts 
        ON workout_sets.workout_id = workouts.id
WHERE workout_sets.exercise_id = $1 
    AND workouts.user_id = $2;

-- name: GetPrVolumeByExerciseIdUserId :one
SELECT COALESCE(
    ROUND(MAX(workout_sets.weight_in_lbs * workout_sets.reps), 2),
    0.0
)::DOUBLE PRECISION as volume
FROM workout_sets 
    INNER JOIN workouts 
        ON workout_sets.workout_id = workouts.id
WHERE workout_sets.exercise_id = $1 
    AND workouts.user_id = $2;