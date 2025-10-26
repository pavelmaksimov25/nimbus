package store

import (
	entity "nimbus/internal/workflow/domain/entity"

	"github.com/google/uuid"
)

type TaskStore interface {
	StoreTask(task *entity.Task) (*entity.Task, error)
	GetTasks() []entity.Task
	GetTask(id uuid.UUID) *entity.Task
}
