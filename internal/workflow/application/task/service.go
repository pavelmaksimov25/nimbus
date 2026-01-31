package task

import (
	"context"
	"log"
	trService "nimbus/internal/task_runnner/domain/service"
	entity "nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/repository"
	"nimbus/internal/workflow/domain/service"
	"nimbus/internal/workflow/domain/types"
	"time"

	"github.com/google/uuid"
)

type taskService struct {
	repository        repository.TaskRepository
	dispatchService   trService.DispatchService
	taskRunnerService trService.TaskRunnerService
}

func NewTaskService(repository repository.TaskRepository, dispatchService trService.DispatchService, taskRunnerService trService.TaskRunnerService) service.TaskService {
	return &taskService{
		repository:        repository,
		dispatchService:   dispatchService,
		taskRunnerService: taskRunnerService,
	}
}

func (ts *taskService) CreateTask(task *entity.Task, runnerID uuid.UUID) (*entity.Task, error) {
	_, err := ts.taskRunnerService.GetRunner(runnerID)
	if err != nil {
		return nil, &types.RecordNotFoundError{Resource: "TaskRunner", ID: runnerID.String()}
	}

	task.ID = uuid.New()
	task.RunnerID = runnerID
	task.Status = entity.StatusNew
	task.CreatedAt = time.Now()

	storedTask, err := ts.repository.StoreTask(task)
	if err != nil {
		return nil, err
	}

	if err := ts.taskRunnerService.AssignTask(runnerID, storedTask.ID); err != nil {
		return nil, err
	}

	return storedTask, nil
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

	if err := ts.repository.UpdateTask(task); err != nil {
		return err
	}

	if ts.dispatchService != nil {
		go func() {
			if err := ts.dispatchService.DispatchTask(context.Background(), task.ID, task.Payload); err != nil {
				log.Printf("failed to dispatch task %s: %v", task.ID, err)
			}
		}()
	}

	return nil
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
