package store

import (
	"sync"

	uuid "github.com/google/uuid"
	entity "nimbus/internal/workflow/domain/entity"
)

type TaskStoreInMemory struct {
	mu    sync.RWMutex
	store map[uuid.UUID]*entity.Task
}

func NewTaskStoreInMemory() TaskStore {
	return &TaskStoreInMemory{
		store: make(map[uuid.UUID]*entity.Task),
	}
}

func (ts *TaskStoreInMemory) StoreTask(task *entity.Task) (*entity.Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.store[task.ID] = task
	return task, nil
}

func (ts *TaskStoreInMemory) GetTasks() []entity.Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	var tasks []entity.Task
	for _, task := range ts.store {
		tasks = append(tasks, *task)
	}
	return tasks
}

func (ts *TaskStoreInMemory) GetTask(id uuid.UUID) *entity.Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task, exists := ts.store[id]
	if !exists {
		return nil
	}
	return task
}

func (ts *TaskStoreInMemory) UpdateTask(task *entity.Task) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.store[task.ID] = task
	return nil
}
