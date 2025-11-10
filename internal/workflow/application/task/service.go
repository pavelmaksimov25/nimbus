package task

import (
	entity "nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/repository"
	"nimbus/internal/workflow/domain/service"
	"nimbus/internal/workflow/domain/types"
	"time"

	"github.com/google/uuid"
)

type taskService struct {
	repository repository.TaskRepository
}

func NewTaskService(repository repository.TaskRepository) service.TaskService {
	return &taskService{
		repository: repository,
	}
}

func (ts *taskService) CreateTask(task *entity.Task) (*entity.Task, error) {
	task.ID = uuid.New()
	task.Status = entity.StatusNew
	task.CreatedAt = time.Now()

	return ts.repository.StoreTask(task)
}

func (ts *taskService) GetTasks() []entity.Task {
	return ts.repository.GetTasks()
}

func (ts *taskService) GetTask(id uuid.UUID) (*entity.Task, error) {
	task := ts.repository.GetTask(id)
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

	return ts.repository.UpdateTask(task)
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

	return ts.repository.UpdateTask(task)
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

	return ts.repository.UpdateTask(task)
}
