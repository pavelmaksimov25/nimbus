package echo

import (
	"context"
	"log"

	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/runner"
)

type echoRunner struct{}

func (r *echoRunner) Execute(_ context.Context, payload string) error {
	log.Printf("[echo runner] payload: %s", payload)
	return nil
}

func NewFactory() runner.Factory {
	return func(_ entity.TaskRunnerConfig) runner.Runner {
		return &echoRunner{}
	}
}
