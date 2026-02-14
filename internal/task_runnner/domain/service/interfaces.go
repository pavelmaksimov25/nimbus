package service

import (
	"context"
	"nimbus/internal/task_runnner/domain/entity"

	"github.com/google/uuid"
)

type TaskRunnerService interface {
	CreateRunner(runner *entity.TaskRunner) (*entity.TaskRunner, error)
	GetRunners() ([]entity.TaskRunner, error)
	GetRunner(id uuid.UUID) (*entity.TaskRunner, error)
	AssignTask(runnerID uuid.UUID, taskID uuid.UUID) error
	UnassignTask(runnerID uuid.UUID, taskID uuid.UUID) error
	GetRunnersByTaskID(taskID uuid.UUID) ([]entity.TaskRunner, error)
}

type DispatchService interface {
	DispatchTask(ctx context.Context, taskID uuid.UUID, payload string) error
}
