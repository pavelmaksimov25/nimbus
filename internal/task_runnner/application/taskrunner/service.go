package taskrunner

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/repository"
	"nimbus/internal/task_runnner/domain/runner"
	"nimbus/internal/task_runnner/domain/service"
	"nimbus/internal/workflow/domain/types"

	"github.com/google/uuid"
)

type taskRunnerService struct {
	repository repository.TaskRunnerRepository
	validators map[entity.TaskRunnerType]runner.ConfigValidator
}

func NewTaskRunnerService(repository repository.TaskRunnerRepository, validators map[entity.TaskRunnerType]runner.ConfigValidator) service.TaskRunnerService {
	return &taskRunnerService{repository: repository, validators: validators}
}

func (s *taskRunnerService) CreateRunner(r *entity.TaskRunner) (*entity.TaskRunner, error) {
	v, ok := s.validators[r.Type]
	if !ok {
		validTypes := make([]string, 0, len(s.validators))
		for t := range s.validators {
			validTypes = append(validTypes, string(t))
		}
		sort.Strings(validTypes)
		return nil, &types.UnprocessableEntityError{
			Msg: fmt.Sprintf("unknown runner type %q, valid types: %s", r.Type, strings.Join(validTypes, ", ")),
		}
	}

	if err := v.Validate(r.Config); err != nil {
		return nil, &types.UnprocessableEntityError{Msg: err.Error()}
	}

	r.ID = uuid.New()
	r.CreatedAt = time.Now()
	return s.repository.Store(r)
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
