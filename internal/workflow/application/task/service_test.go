package task

import (
	"errors"
	"testing"
	"time"

	trEntity "nimbus/internal/task_runnner/domain/entity"
	trMocks "nimbus/internal/task_runnner/application/mocks"
	"nimbus/internal/workflow/adapters/mocks"
	"nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/types"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTaskService_CreateTask(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	mockRunnerSvc := trMocks.NewMockTaskRunnerService(ctrl)
	taskService := NewTaskService(mockRepo, nil, mockRunnerSvc)

	runnerID := uuid.New()
	inputTask := &entity.Task{
		Payload:    "Test Payload",
		WorkflowID: uuid.New(),
	}

	mockRunnerSvc.EXPECT().
		GetRunner(runnerID).
		Return(&trEntity.TaskRunner{ID: runnerID, Name: "Echo Runner", Type: trEntity.Echo}, nil)

	mockRepo.EXPECT().
		StoreTask(gomock.Any()).
		DoAndReturn(func(task *entity.Task) (*entity.Task, error) {
			return task, nil
		})

	mockRunnerSvc.EXPECT().
		AssignTask(runnerID, gomock.Any()).
		Return(nil)

	// Act
	resultTask, err := taskService.CreateTask(inputTask, runnerID)

	// Assert
	assert.NoError(t, err)
	assert.IsType(t, uuid.UUID{}, resultTask.ID)
	assert.Equal(t, "Test Payload", resultTask.Payload)
	assert.Equal(t, entity.StatusNew, resultTask.Status)
	assert.Equal(t, runnerID, resultTask.RunnerID)
}

func TestTaskService_CreateTask_RunnerNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	mockRunnerSvc := trMocks.NewMockTaskRunnerService(ctrl)
	taskService := NewTaskService(mockRepo, nil, mockRunnerSvc)

	runnerID := uuid.New()
	inputTask := &entity.Task{
		Payload:    "Test Payload",
		WorkflowID: uuid.New(),
	}

	mockRunnerSvc.EXPECT().
		GetRunner(runnerID).
		Return(nil, &types.RecordNotFoundError{Resource: "TaskRunner", ID: runnerID.String()})

	// Act
	resultTask, err := taskService.CreateTask(inputTask, runnerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultTask)
	assert.IsType(t, &types.RecordNotFoundError{}, err)
}

func TestTaskService_GetTasks(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo, nil, nil)

	expectedTask := entity.Task{
		ID:         uuid.New(),
		Payload:    "Test Payload",
		Status:     entity.StatusNew,
		WorkflowID: uuid.New(),
	}

	mockRepo.EXPECT().
		GetTasks().
		Return([]entity.Task{expectedTask})

	// Act
	tasks := taskService.GetTasks()

	// Assert
	assert.Len(t, tasks, 1)
	assert.Equal(t, expectedTask.Payload, tasks[0].Payload)
	assert.Equal(t, expectedTask.Status, tasks[0].Status)
	assert.Equal(t, expectedTask.WorkflowID, tasks[0].WorkflowID)
}

func TestTaskService_GetTask(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo, nil, nil)

	expectedTask := &entity.Task{
		ID:         uuid.New(),
		Payload:    "Test Payload",
		Status:     entity.StatusNew,
		WorkflowID: uuid.New(),
	}

	mockRepo.EXPECT().
		GetTask(expectedTask.ID).
		Return(expectedTask)

	// Act
	retrievedTask, err := taskService.GetTask(expectedTask.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTask.Payload, retrievedTask.Payload)
	assert.Equal(t, expectedTask.Status, retrievedTask.Status)
	assert.Equal(t, expectedTask.WorkflowID, retrievedTask.WorkflowID)
}

func TestTaskService_GetTask_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo, nil, nil)
	nonExistentID := uuid.New()

	mockRepo.EXPECT().
		GetTask(nonExistentID).
		Return(nil)

	// Act
	retrievedTask, err := taskService.GetTask(nonExistentID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrievedTask)
}

func TestTaskService_StartTask(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo, nil, nil)

	existingTask := &entity.Task{
		ID:         uuid.New(),
		Payload:    "Test Payload",
		Status:     entity.StatusNew,
		WorkflowID: uuid.New(),
	}

	mockRepo.EXPECT().
		GetTask(existingTask.ID).
		Return(existingTask)

	mockRepo.EXPECT().
		UpdateTaskStatus(existingTask.ID, entity.StatusNew, entity.StatusInProgress).
		Return(nil)

	// Act
	err := taskService.StartTask(existingTask.ID)

	// Assert
	assert.Nil(t, err)
}

func TestTaskService_StartTask_InvalidTaskStatus(t *testing.T) {
	// Arrange
	tests := []entity.Task{
		{
			ID:         uuid.New(),
			Payload:    "test status in_progress",
			Status:     entity.StatusInProgress,
			WorkflowID: uuid.New(),
		},
		{
			ID:         uuid.New(),
			Payload:    "test status failed",
			Status:     entity.StatusFailed,
			WorkflowID: uuid.New(),
		},
		{
			ID:         uuid.New(),
			Payload:    "test status completed",
			Status:     entity.StatusCompleted,
			WorkflowID: uuid.New(),
		},
	}

	for _, existingTask := range tests {
		t.Run(existingTask.Payload, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTaskRepository(ctrl)
			taskService := NewTaskService(mockRepo, nil, nil)

			taskCopy := existingTask
			mockRepo.EXPECT().
				GetTask(existingTask.ID).
				Return(&taskCopy)

			mockRepo.EXPECT().
				UpdateTaskStatus(existingTask.ID, entity.StatusNew, entity.StatusInProgress).
				Return(errors.New("task is not in status NEW"))

			// Act
			err := taskService.StartTask(existingTask.ID)

			// Assert
			assert.IsType(t, &types.UnprocessableEntityError{}, err)
		})
	}
}

func TestTaskService_StartTask_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo, nil, nil)
	nonExistentID := uuid.New()

	mockRepo.EXPECT().
		GetTask(nonExistentID).
		Return(nil)

	// Act
	err := taskService.StartTask(nonExistentID)

	// Assert
	assert.IsType(t, &types.RecordNotFoundError{}, err)
}

func TestTaskService_CompleteTask(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo, nil, nil)

	taskID := uuid.New()
	inProgressTask := &entity.Task{
		ID:         taskID,
		Payload:    "Test Payload",
		Status:     entity.StatusInProgress,
		WorkflowID: uuid.New(),
	}

	// First GetTask call for existence check
	mockRepo.EXPECT().
		GetTask(taskID).
		Return(inProgressTask)

	mockRepo.EXPECT().
		UpdateTaskStatus(taskID, entity.StatusInProgress, entity.StatusCompleted).
		Return(nil)

	// Second GetTask call for payload update
	completedTask := &entity.Task{
		ID:         taskID,
		Payload:    "Test Payload",
		Status:     entity.StatusCompleted,
		WorkflowID: inProgressTask.WorkflowID,
	}
	mockRepo.EXPECT().
		GetTask(taskID).
		Return(completedTask)

	mockRepo.EXPECT().
		UpdateTask(gomock.Any()).
		DoAndReturn(func(task *entity.Task) error {
			assert.Equal(t, "Test Payload - Additional Payload", task.Payload)
			return nil
		})

	// Act
	err := taskService.CompleteTask(taskID, " - Additional Payload")

	// Assert
	assert.Nil(t, err)
}

func TestTaskService_CompleteTask_InvalidTaskStatus(t *testing.T) {
	// Arrange
	tests := []entity.Task{
		{
			ID:         uuid.New(),
			Payload:    "test status new",
			Status:     entity.StatusNew,
			WorkflowID: uuid.New(),
		},
		{
			ID:         uuid.New(),
			Payload:    "test status failed",
			Status:     entity.StatusFailed,
			WorkflowID: uuid.New(),
		},
		{
			ID:         uuid.New(),
			Payload:    "test status completed",
			Status:     entity.StatusCompleted,
			WorkflowID: uuid.New(),
		},
	}

	for _, existingTask := range tests {
		t.Run(existingTask.Payload, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTaskRepository(ctrl)
			taskService := NewTaskService(mockRepo, nil, nil)

			taskCopy := existingTask
			mockRepo.EXPECT().
				GetTask(existingTask.ID).
				Return(&taskCopy)

			mockRepo.EXPECT().
				UpdateTaskStatus(existingTask.ID, entity.StatusInProgress, entity.StatusCompleted).
				Return(errors.New("task is not in status IN_PROGRESS"))

			// Act
			err := taskService.CompleteTask(existingTask.ID, "")

			// Assert
			assert.IsType(t, &types.UnprocessableEntityError{}, err)
		})
	}
}

func TestTaskService_CompleteTask_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo, nil, nil)
	nonExistentID := uuid.New()

	mockRepo.EXPECT().
		GetTask(nonExistentID).
		Return(nil)

	// Act
	err := taskService.CompleteTask(nonExistentID, "")

	// Assert
	assert.IsType(t, &types.RecordNotFoundError{}, err)
}

func TestTaskService_FailTask(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo, nil, nil)

	taskID := uuid.New()
	inProgressTask := &entity.Task{
		ID:      taskID,
		Payload: "Test Payload",
		Status:  entity.StatusInProgress,
	}

	mockRepo.EXPECT().
		GetTask(taskID).
		Return(inProgressTask)

	mockRepo.EXPECT().
		UpdateTaskStatus(taskID, entity.StatusInProgress, entity.StatusFailed).
		Return(nil)

	failedTask := &entity.Task{
		ID:      taskID,
		Payload: "Test Payload",
		Status:  entity.StatusFailed,
	}
	mockRepo.EXPECT().
		GetTask(taskID).
		Return(failedTask)

	mockRepo.EXPECT().
		UpdateTask(gomock.Any()).
		DoAndReturn(func(task *entity.Task) error {
			assert.Equal(t, "Some failure reason", task.FailReason)
			return nil
		})

	// Act
	err := taskService.FailTask(taskID, "Some failure reason")

	// Assert
	assert.Nil(t, err)
}

func TestTaskService_FailTask_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo, nil, nil)
	nonExistentID := uuid.New()

	mockRepo.EXPECT().
		GetTask(nonExistentID).
		Return(nil)

	// Act
	err := taskService.FailTask(nonExistentID, "Some failure reason")

	// Assert
	assert.IsType(t, &types.RecordNotFoundError{}, err)
}

func TestTaskService_FailTask_InvalidTaskStatus(t *testing.T) {
	// Arrange
	tests := []entity.Task{
		{
			ID:      uuid.New(),
			Payload: "test status new",
			Status:  entity.StatusNew,
		},
		{
			ID:      uuid.New(),
			Payload: "test status completed",
			Status:  entity.StatusCompleted,
		},
		{
			ID:      uuid.New(),
			Payload: "test status failed",
			Status:  entity.StatusFailed,
		},
	}

	for _, existingTask := range tests {
		t.Run(existingTask.Payload, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTaskRepository(ctrl)
			taskService := NewTaskService(mockRepo, nil, nil)

			taskCopy := existingTask
			mockRepo.EXPECT().
				GetTask(existingTask.ID).
				Return(&taskCopy)

			mockRepo.EXPECT().
				UpdateTaskStatus(existingTask.ID, entity.StatusInProgress, entity.StatusFailed).
				Return(errors.New("task is not in status IN_PROGRESS"))

			// Act
			err := taskService.FailTask(existingTask.ID, "Some failure reason")

			// Assert
			assert.IsType(t, &types.UnprocessableEntityError{}, err)
		})
	}
}

func TestTaskService_StartTask_DispatcheToRunners(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	mockDispatch := trMocks.NewMockDispatchService(ctrl)
	taskService := NewTaskService(mockRepo, mockDispatch, nil)

	existingTask := &entity.Task{
		ID:         uuid.New(),
		Payload:    "Test Payload",
		Status:     entity.StatusNew,
		WorkflowID: uuid.New(),
	}

	mockRepo.EXPECT().
		GetTask(existingTask.ID).
		Return(existingTask)

	mockRepo.EXPECT().
		UpdateTaskStatus(existingTask.ID, entity.StatusNew, entity.StatusInProgress).
		Return(nil)

	mockDispatch.EXPECT().
		DispatchTask(gomock.Any(), existingTask.ID, existingTask.Payload).
		Return(nil)

	// Act
	err := taskService.StartTask(existingTask.ID)

	// Assert
	assert.Nil(t, err)

	// Wait for the goroutine to complete
	time.Sleep(50 * time.Millisecond)
}
