package dispatch

import (
	"context"
	"errors"
	"testing"

	"nimbus/internal/task_runnner/adapters/mocks"
	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/runner"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockRunner struct {
	err error
}

func (r *mockRunner) Execute(ctx context.Context, payload string) error {
	return r.err
}

func TestDispatchService_DispatchTask_FanOut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)

	executedCount := 0
	factories := map[entity.TaskRunnerType]runner.Factory{
		entity.AwsSqs: func(config entity.TaskRunnerConfig) runner.Runner {
			return &mockRunner{err: nil}
		},
	}

	svc := NewDispatchService(mockRepo, factories)

	taskID := uuid.New()
	runners := []entity.TaskRunner{
		{ID: uuid.New(), Name: "runner-1", Type: entity.AwsSqs, Config: entity.TaskRunnerConfig{"queue_url": "url1"}},
		{ID: uuid.New(), Name: "runner-2", Type: entity.AwsSqs, Config: entity.TaskRunnerConfig{"queue_url": "url2"}},
	}

	mockRepo.EXPECT().GetByTaskID(taskID).Return(runners, nil)

	// Override factory to count executions
	factories[entity.AwsSqs] = func(config entity.TaskRunnerConfig) runner.Runner {
		executedCount++
		return &mockRunner{err: nil}
	}

	err := svc.DispatchTask(context.Background(), taskID, "test payload")

	assert.NoError(t, err)
	assert.Equal(t, 2, executedCount)
}

func TestDispatchService_DispatchTask_NoRunners(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)

	factories := map[entity.TaskRunnerType]runner.Factory{}
	svc := NewDispatchService(mockRepo, factories)

	taskID := uuid.New()
	mockRepo.EXPECT().GetByTaskID(taskID).Return([]entity.TaskRunner{}, nil)

	err := svc.DispatchTask(context.Background(), taskID, "test payload")

	assert.NoError(t, err)
}

func TestDispatchService_DispatchTask_RunnerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)

	factories := map[entity.TaskRunnerType]runner.Factory{
		entity.AwsSqs: func(config entity.TaskRunnerConfig) runner.Runner {
			return &mockRunner{err: errors.New("send failed")}
		},
	}

	svc := NewDispatchService(mockRepo, factories)

	taskID := uuid.New()
	runners := []entity.TaskRunner{
		{ID: uuid.New(), Name: "runner-1", Type: entity.AwsSqs, Config: entity.TaskRunnerConfig{"queue_url": "url1"}},
	}

	mockRepo.EXPECT().GetByTaskID(taskID).Return(runners, nil)

	err := svc.DispatchTask(context.Background(), taskID, "test payload")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "send failed")
}

func TestDispatchService_DispatchTask_UnknownType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)

	factories := map[entity.TaskRunnerType]runner.Factory{}
	svc := NewDispatchService(mockRepo, factories)

	taskID := uuid.New()
	runners := []entity.TaskRunner{
		{ID: uuid.New(), Name: "runner-1", Type: "unknown", Config: entity.TaskRunnerConfig{}},
	}

	mockRepo.EXPECT().GetByTaskID(taskID).Return(runners, nil)

	err := svc.DispatchTask(context.Background(), taskID, "test payload")

	assert.NoError(t, err)
}

func TestDispatchService_DispatchTask_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)

	factories := map[entity.TaskRunnerType]runner.Factory{}
	svc := NewDispatchService(mockRepo, factories)

	taskID := uuid.New()
	mockRepo.EXPECT().GetByTaskID(taskID).Return(nil, errors.New("db error"))

	err := svc.DispatchTask(context.Background(), taskID, "test payload")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get runners")
}

func TestDispatchService_DispatchTask_MultipleErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)

	callCount := 0
	factories := map[entity.TaskRunnerType]runner.Factory{
		entity.AwsSqs: func(config entity.TaskRunnerConfig) runner.Runner {
			callCount++
			return &mockRunner{err: errors.New("fail-" + config["queue_url"].(string))}
		},
	}

	svc := NewDispatchService(mockRepo, factories)

	taskID := uuid.New()
	runners := []entity.TaskRunner{
		{ID: uuid.New(), Name: "runner-1", Type: entity.AwsSqs, Config: entity.TaskRunnerConfig{"queue_url": "url1"}},
		{ID: uuid.New(), Name: "runner-2", Type: entity.AwsSqs, Config: entity.TaskRunnerConfig{"queue_url": "url2"}},
	}

	mockRepo.EXPECT().GetByTaskID(taskID).Return(runners, nil)

	err := svc.DispatchTask(context.Background(), taskID, "test payload")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fail-url1")
	assert.Contains(t, err.Error(), "fail-url2")
	assert.Equal(t, 2, callCount)
}
