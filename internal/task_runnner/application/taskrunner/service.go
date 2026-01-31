package taskrunner

import (
	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/repository"
	"nimbus/internal/task_runnner/domain/service"
	"nimbus/internal/workflow/domain/types"
	"time"

	"github.com/google/uuid"
)

type taskRunnerService struct {
	repository repository.TaskRunnerRepository
}

func NewTaskRunnerService(repository repository.TaskRunnerRepository) service.TaskRunnerService {
	return &taskRunnerService{repository: repository}
}

func (s *taskRunnerService) CreateRunner(runner *entity.TaskRunner) (*entity.TaskRunner, error) {
	runner.ID = uuid.New()
	runner.CreatedAt = time.Now()
	return s.repository.Store(runner)
}

func (s *taskRunnerService) GetRunners() ([]entity.TaskRunner, error) {
	return s.repository.GetAll()
}

func (s *taskRunnerService) GetRunner(id uuid.UUID) (*entity.TaskRunner, error) {
	runner, err := s.repository.GetByID(id)
	if err != nil {
		return nil, &types.RecordNotFoundError{Resource: "TaskRunner", ID: id.String()}
	}
	return runner, nil
}

func (s *taskRunnerService) AssignTask(runnerID uuid.UUID, taskID uuid.UUID) error {
	_, err := s.repository.GetByID(runnerID)
	if err != nil {
		return &types.RecordNotFoundError{Resource: "TaskRunner", ID: runnerID.String()}
	}
	return s.repository.AssignTask(runnerID, taskID)
}

func (s *taskRunnerService) UnassignTask(runnerID uuid.UUID, taskID uuid.UUID) error {
	return s.repository.UnassignTask(runnerID, taskID)
}

func (s *taskRunnerService) GetRunnersByTaskID(taskID uuid.UUID) ([]entity.TaskRunner, error) {
	return s.repository.GetByTaskID(taskID)
}
