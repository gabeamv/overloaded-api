-- name: CreateExercise :one
INSERT INTO exercises (name, is_custom, user_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetAllExercisesUserId :many
SELECT * FROM exercises
WHERE user_id = $1 AND is_custom = true;

-- name: GetExerciseById :one
SELECT * FROM exercises
WHERE id = $1;

-- name: UpdateExerciseNameById :one
UPDATE exercises
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteExerciseById :exec
DELETE FROM exercises
WHERE id = $1;