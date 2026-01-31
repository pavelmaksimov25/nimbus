-- +goose Up
CREATE TYPE task_runner_type AS ENUM ('queue');

CREATE TABLE task_runners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type task_runner_type NOT NULL DEFAULT 'queue',
    config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE task_runners_tasks (
    task_runner_id UUID NOT NULL REFERENCES task_runners(id) ON DELETE CASCADE,
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    PRIMARY KEY (task_runner_id, task_id)
);

-- +goose Down
DROP TABLE IF EXISTS task_runners_tasks;
DROP TABLE IF EXISTS task_runners;
DROP TYPE IF EXISTS task_runner_type;
