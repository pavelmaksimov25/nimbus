package task

import (
	store "nimbus/internal/workflow/adapters/store"
	entity "nimbus/internal/workflow/domain/entity"
)


type Service interface {
	CreateTask(payload string) (string, error)
	GetTasks() ([]entity.Task)
}

type taskService struct {
	store store.TaskStore
}

func NewTaskService(store store.TaskStore) Service {
	return &taskService{
		store: store,
	}
}

func (ts *taskService) CreateTask(payload string) (string, error) {
	id, err := ts.store.StoreTask(payload)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (ts *taskService) GetTasks() ([]entity.Task) {
	return ts.store.GetTasks()
}
