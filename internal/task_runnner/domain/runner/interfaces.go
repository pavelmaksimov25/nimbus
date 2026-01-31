package runner

import (
	"context"
	"nimbus/internal/task_runnner/domain/entity"
)

type Runner interface {
	Execute(ctx context.Context, payload string) error
}

type Factory func(config entity.TaskRunnerConfig) Runner

type ConfigValidator interface {
	Validate(config entity.TaskRunnerConfig) error
}
