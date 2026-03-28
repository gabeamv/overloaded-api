-- +goose Up
CREATE TABLE workout_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_id UUID NOT NULL,
    FOREIGN KEY (exercise_id) REFERENCES exercises (id),
    weight_in_lbs NUMERIC(6,2),
    reps NUMERIC(6,2),
    time_in_seconds numeric(10,2),
    progress_track UUID NOT NULL,
    FOREIGN KEY (progress_track) REFERENCES progression_rules(id)
);

-- +goose Down
DROP TABLE workout_sets;