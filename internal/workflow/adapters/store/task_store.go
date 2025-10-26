package store

import (
	"sync"
	"time"

	entity "nimbus/internal/workflow/domain/entity"

	uuid "github.com/google/uuid"
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

func (ts *TaskStoreInMemory) StoreTask(payload string) (*entity.Task, error) {
	id := uuid.New()
	task := &entity.Task{
		ID:      id,
		Payload: payload,
		CreatedAt: time.Now(),
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.store[id] = task

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
