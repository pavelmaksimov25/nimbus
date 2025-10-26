package task

import (
	store "nimbus/internal/workflow/adapters/store"
	entity "nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/types"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateTask(payload string) (*entity.Task, error)
	GetTasks() []entity.Task
	GetTask(id uuid.UUID) (*entity.Task, error)
}

type taskService struct {
	store store.TaskStore
}

func NewTaskService(store store.TaskStore) Service {
	return &taskService{
		store: store,
	}
}

func (ts *taskService) CreateTask(payload string) (*entity.Task, error) {
	id := uuid.New()
	task := &entity.Task{
		ID:        id,
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	return ts.store.StoreTask(task)
}

func (ts *taskService) GetTasks() []entity.Task {
	return ts.store.GetTasks()
}

func (ts *taskService) GetTask(id uuid.UUID) (*entity.Task, error) {
	task := ts.store.GetTask(id)
	if task == nil {
		return nil, &types.RecordNotFoundError{Resource: "Task", ID: id.String()}
	}
	return task, nil
}
