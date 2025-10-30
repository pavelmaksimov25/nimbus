package task

import (
	store "nimbus/internal/workflow/adapters/storage"
	entity "nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/types"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateTask(task *entity.Task) (*entity.Task, error)
	GetTasks() []entity.Task
	GetTask(id uuid.UUID) (*entity.Task, error)
	StartTask(id uuid.UUID) error
	CompleteTask(id uuid.UUID, additionalPayload string) error
	FailTask(id uuid.UUID, reason string) error
}

type taskService struct {
	store store.TaskStorage
}

func NewTaskService(store store.TaskStorage) Service {
	return &taskService{
		store: store,
	}
}

func (ts *taskService) CreateTask(task *entity.Task) (*entity.Task, error) {
	task.ID = uuid.New()
	task.Status = entity.StatusNew
	task.CreatedAt = time.Now()

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

func (ts *taskService) StartTask(id uuid.UUID) error {
	task, err := ts.GetTask(id)
	if err != nil {
		return err
	}

	if task.Status != entity.StatusNew {
		return &types.UnprocessableEntityError{
			Resource: "Task",
			ID:       id.String(),
			Msg:      "only tasks with status 'new' can be started",
		}
	}

	task.Status = entity.StatusInProgress

	return ts.store.UpdateTask(task)
}

func (ts *taskService) CompleteTask(id uuid.UUID, additionalPayload string) error {
	task, err := ts.GetTask(id)
	if err != nil {
		return err
	}

	if task.Status != entity.StatusInProgress {
		return &types.UnprocessableEntityError{
			Resource: "Task",
			ID:       id.String(),
			Msg:      "only tasks with status 'in_progress' can be completed",
		}
	}

	task.Status = entity.StatusCompleted
	task.Payload += additionalPayload

	return ts.store.UpdateTask(task)
}

func (ts *taskService) FailTask(id uuid.UUID, reason string) error {
	task, err := ts.GetTask(id)
	if err != nil {
		return err
	}

	// it is unable to fail what wasn't started
	// it is also unable to fail what is already failed or completed
	if task.Status != entity.StatusInProgress {
		return &types.UnprocessableEntityError{
			Resource: "Task",
			ID:       id.String(),
			Msg:      "only tasks with status 'in_progress' can be failed",
		}
	}

	task.Status = entity.StatusFailed
	task.FailReason = reason

	return ts.store.UpdateTask(task)
}
