package service

import (
	"nimbus/internal/workflow/domain/entity"

	"github.com/google/uuid"
)


type WorkflowService interface {
	CreateWorkflow(name string) (*entity.Workflow, error)
}

type TaskService interface {
	CreateTask(task *entity.Task) (*entity.Task, error)
	GetTasks() []entity.Task
	GetTask(id uuid.UUID) (*entity.Task, error)
	StartTask(id uuid.UUID) error
	CompleteTask(id uuid.UUID, additionalPayload string) error
	FailTask(id uuid.UUID, reason string) error
}