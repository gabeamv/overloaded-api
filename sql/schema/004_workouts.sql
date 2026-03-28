-- +goose Up
CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP NOT NULL,
    volume INTEGER NOT NULL,
    prs INTEGER NOT NULL
);

-- +goose Down
DROP TABLE workouts;