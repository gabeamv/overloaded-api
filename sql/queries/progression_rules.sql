-- name: CreateProgressionRule :exec
INSERT INTO progression_rules (label, rule, description)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetAllProgressionRules :many
SELECT * FROM progression_rules;

-- name: DeleteAllProgressionRules :exec
DELETE FROM progression_rules;