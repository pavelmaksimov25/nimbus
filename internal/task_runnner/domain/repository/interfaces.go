package repository

import (
	"nimbus/internal/task_runnner/domain/entity"

	"github.com/google/uuid"
)

type TaskRunnerRepository interface {
	Store(runner *entity.TaskRunner) (*entity.TaskRunner, error)
	GetAll() ([]entity.TaskRunner, error)
	GetByID(id uuid.UUID) (*entity.TaskRunner, error)
	AssignTask(runnerID uuid.UUID, taskID uuid.UUID) error
	UnassignTask(runnerID uuid.UUID, taskID uuid.UUID) error
	GetByTaskID(taskID uuid.UUID) ([]entity.TaskRunner, error)
}
