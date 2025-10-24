package store

import (
	entity "nimbus/internal/workflow/domain/entity"
	uuid "github.com/google/uuid"
)

type TaskStore interface {
	StoreTask(payload string) (uuid.UUID, error)
	GetTasks() []entity.Task
}