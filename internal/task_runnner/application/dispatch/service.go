package dispatch

import (
	"context"
	"fmt"
	"log"
	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/repository"
	"nimbus/internal/task_runnner/domain/runner"
	"nimbus/internal/task_runnner/domain/service"

	"github.com/google/uuid"
)

type dispatchService struct {
	repository repository.TaskRunnerRepository
	factories  map[entity.TaskRunnerType]runner.Factory
}

func NewDispatchService(
	repository repository.TaskRunnerRepository,
	factories map[entity.TaskRunnerType]runner.Factory,
) service.DispatchService {
	return &dispatchService{
		repository: repository,
		factories:  factories,
	}
}

func (s *dispatchService) DispatchTask(ctx context.Context, taskID uuid.UUID, payload string) error {
	runners, err := s.repository.GetByTaskID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get runners for task %s: %w", taskID, err)
	}

	if len(runners) == 0 {
		return nil
	}

	var dispatchErr error
	for _, r := range runners {
		factory, ok := s.factories[r.Type]
		if !ok {
			log.Printf("no factory registered for runner type %q, skipping runner %s", r.Type, r.ID)
			continue
		}

		execRunner := factory(r.Config)
		if err := execRunner.Execute(ctx, payload); err != nil {
			log.Printf("failed to dispatch task %s to runner %s: %v", taskID, r.ID, err)
			dispatchErr = err
		}
	}

	return dispatchErr
}
