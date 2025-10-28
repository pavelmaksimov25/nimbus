package storage

import (
	"sync"

	entity "nimbus/internal/workflow/domain/entity"

	uuid "github.com/google/uuid"
)

type TaskStorageInMemory struct {
	mu    sync.RWMutex
	store map[uuid.UUID]*entity.Task
}

func NewTaskStorageInMemory() TaskStorage {
	return &TaskStorageInMemory{
		store: make(map[uuid.UUID]*entity.Task),
	}
}

func (ts *TaskStorageInMemory) StoreTask(task *entity.Task) (*entity.Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.store[task.ID] = task
	return task, nil
}

func (ts *TaskStorageInMemory) GetTasks() []entity.Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	var tasks []entity.Task
	for _, task := range ts.store {
		tasks = append(tasks, *task)
	}
	return tasks
}

func (ts *TaskStorageInMemory) GetTask(id uuid.UUID) *entity.Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task, exists := ts.store[id]
	if !exists {
		return nil
	}
	return task
}

func (ts *TaskStorageInMemory) UpdateTask(task *entity.Task) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.store[task.ID] = task
	return nil
}
