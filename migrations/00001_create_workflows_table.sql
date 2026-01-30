-- +goose Up
CREATE TABLE IF NOT EXISTS workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_workflows_created_at ON workflows(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_workflows_created_at;
DROP TABLE IF EXISTS workflows;
