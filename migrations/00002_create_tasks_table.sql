-- +goose Up
CREATE TYPE task_status AS ENUM ('NEW', 'IN_PROGRESS', 'COMPLETED', 'FAILED');

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payload TEXT NOT NULL,
    workflow_id UUID NOT NULL,
    status task_status NOT NULL DEFAULT 'NEW',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fail_reason TEXT,
    CONSTRAINT fk_workflow
        FOREIGN KEY(workflow_id)
        REFERENCES workflows(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_tasks_workflow_id ON tasks(workflow_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_tasks_created_at;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_workflow_id;
DROP TABLE IF EXISTS tasks;
DROP TYPE IF EXISTS task_status;
