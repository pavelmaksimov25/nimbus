package echo

import (
	"context"
	"log"

	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/runner"
)

type echoRunner struct{}

const maxLogPayloadLen = 256

func (r *echoRunner) Execute(_ context.Context, payload string) error {
	logPayload := payload
	if len(logPayload) > maxLogPayloadLen {
		logPayload = logPayload[:maxLogPayloadLen] + "...(truncated)"
	}
	log.Printf("[echo runner] payload: %s", logPayload)
	return nil
}

func NewFactory() runner.Factory {
	return func(_ entity.TaskRunnerConfig) runner.Runner {
		return &echoRunner{}
	}
}
