-- +goose Up
ALTER TABLE tasks ADD COLUMN runner_id UUID NOT NULL REFERENCES task_runners(id);
CREATE INDEX idx_tasks_runner_id ON tasks(runner_id);

-- +goose Down
DROP INDEX IF EXISTS idx_tasks_runner_id;
ALTER TABLE tasks DROP COLUMN runner_id;
