package store

import (
	entity "nimbus/internal/workflow/domain/entity"
)

type TaskStore interface {
	StoreTask(payload string) (*entity.Task, error)
	GetTasks() []entity.Task
}
