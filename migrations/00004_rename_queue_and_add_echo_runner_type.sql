-- +goose Up
ALTER TYPE task_runner_type RENAME VALUE 'queue' TO 'aws_sqs';
ALTER TYPE task_runner_type ADD VALUE 'echo';
ALTER TABLE task_runners ALTER COLUMN type SET DEFAULT 'aws_sqs';

-- +goose Down
ALTER TABLE task_runners ALTER COLUMN type SET DEFAULT 'queue';
-- Note: PostgreSQL does not support removing or renaming enum values directly.
-- A full enum rebuild would be needed for a true rollback.
