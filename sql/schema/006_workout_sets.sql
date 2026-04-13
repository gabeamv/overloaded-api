-- +goose Up
CREATE TABLE workout_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    progress_track UUID NOT NULL REFERENCES progression_rules(id),
    weight_in_lbs NUMERIC(6,2) NOT NULL DEFAULT 0,
    reps NUMERIC(6,2) NOT NULL DEFAULT 0,
    time_in_seconds numeric(10,2) NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE workout_sets;