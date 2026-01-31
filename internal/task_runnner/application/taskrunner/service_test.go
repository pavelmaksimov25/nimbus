package taskrunner

import (
	"testing"

	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/repository/mocks"
	"nimbus/internal/workflow/domain/types"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestTaskRunnerService_CreateRunner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)
	svc := NewTaskRunnerService(mockRepo)

	input := &entity.TaskRunner{
		Name:   "my-sqs-runner",
		Type:   entity.Queue,
		Config: entity.TaskRunnerConfig{"queue_url": "http://localhost:9324/000000000000/nimbus-tasks"},
	}

	mockRepo.EXPECT().
		Store(gomock.Any()).
		DoAndReturn(func(runner *entity.TaskRunner) (*entity.TaskRunner, error) {
			return runner, nil
		})

	result, err := svc.CreateRunner(input)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, "my-sqs-runner", result.Name)
	assert.Equal(t, entity.Queue, result.Type)
	assert.False(t, result.CreatedAt.IsZero())
}

func TestTaskRunnerService_GetRunners(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)
	svc := NewTaskRunnerService(mockRepo)

	expected := []entity.TaskRunner{
		{ID: uuid.New(), Name: "runner-1", Type: entity.Queue},
		{ID: uuid.New(), Name: "runner-2", Type: entity.Queue},
	}

	mockRepo.EXPECT().GetAll().Return(expected, nil)

	runners, err := svc.GetRunners()

	assert.NoError(t, err)
	assert.Len(t, runners, 2)
}

func TestTaskRunnerService_GetRunner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)
	svc := NewTaskRunnerService(mockRepo)

	expected := &entity.TaskRunner{ID: uuid.New(), Name: "runner-1", Type: entity.Queue}

	mockRepo.EXPECT().GetByID(expected.ID).Return(expected, nil)

	runner, err := svc.GetRunner(expected.ID)

	assert.NoError(t, err)
	assert.Equal(t, expected.Name, runner.Name)
}

func TestTaskRunnerService_GetRunner_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)
	svc := NewTaskRunnerService(mockRepo)

	id := uuid.New()
	mockRepo.EXPECT().GetByID(id).Return(nil, gorm.ErrRecordNotFound)

	runner, err := svc.GetRunner(id)

	assert.Nil(t, runner)
	assert.IsType(t, &types.RecordNotFoundError{}, err)
}

func TestTaskRunnerService_AssignTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)
	svc := NewTaskRunnerService(mockRepo)

	runnerID := uuid.New()
	taskID := uuid.New()

	mockRepo.EXPECT().
		GetByID(runnerID).
		Return(&entity.TaskRunner{ID: runnerID}, nil)

	mockRepo.EXPECT().
		AssignTask(runnerID, taskID).
		Return(nil)

	err := svc.AssignTask(runnerID, taskID)

	assert.NoError(t, err)
}

func TestTaskRunnerService_AssignTask_RunnerNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)
	svc := NewTaskRunnerService(mockRepo)

	runnerID := uuid.New()
	taskID := uuid.New()

	mockRepo.EXPECT().
		GetByID(runnerID).
		Return(nil, gorm.ErrRecordNotFound)

	err := svc.AssignTask(runnerID, taskID)

	assert.IsType(t, &types.RecordNotFoundError{}, err)
}

func TestTaskRunnerService_UnassignTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)
	svc := NewTaskRunnerService(mockRepo)

	runnerID := uuid.New()
	taskID := uuid.New()

	mockRepo.EXPECT().
		UnassignTask(runnerID, taskID).
		Return(nil)

	err := svc.UnassignTask(runnerID, taskID)

	assert.NoError(t, err)
}

func TestTaskRunnerService_GetRunnersByTaskID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRunnerRepository(ctrl)
	svc := NewTaskRunnerService(mockRepo)

	taskID := uuid.New()
	expected := []entity.TaskRunner{
		{ID: uuid.New(), Name: "runner-1", Type: entity.Queue},
	}

	mockRepo.EXPECT().GetByTaskID(taskID).Return(expected, nil)

	runners, err := svc.GetRunnersByTaskID(taskID)

	assert.NoError(t, err)
	assert.Len(t, runners, 1)
}
