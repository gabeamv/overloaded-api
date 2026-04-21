-- +goose Up
CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP NOT NULL,
    volume INTEGER NOT NULL DEFAULT 0,
    prs INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE workouts;