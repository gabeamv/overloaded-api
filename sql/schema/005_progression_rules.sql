-- +goose Up
CREATE TABLE progression_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    label CHAR(1) UNIQUE NOT NULL,
    rule TEXT NOT NULL,
    description TEXT NOT NULL
);

-- +goose Down
DROP TABLE progression_rules;