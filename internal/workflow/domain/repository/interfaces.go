package repository

import (
	"nimbus/internal/workflow/domain/entity"

	"github.com/google/uuid"
)

type TaskRepository interface {
	StoreTask(task *entity.Task) (*entity.Task, error)
	GetTasks() []entity.Task
	GetTask(id uuid.UUID) *entity.Task
	UpdateTask(task *entity.Task) error
}

type WorkflowRepository interface {
	StoreWorkflow(workflow *entity.Workflow) (*entity.Workflow, error)
	GetWorkflows() []entity.Workflow
	GetWorkflow(id uuid.UUID) *entity.Workflow
	UpdateWorkflow(workflow *entity.Workflow) error
}
