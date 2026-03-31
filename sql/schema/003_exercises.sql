-- +goose Up
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    is_custom BOOLEAN NOT NULL DEFAULT FALSE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    CHECK (
        (is_custom = true AND user_id IS NOT NULL)
        OR
        (is_custom = false AND user_id IS NULL)
    )
);

-- +goose Down
DROP TABLE exercises;