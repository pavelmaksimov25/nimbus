package task

import (
	"testing"

	store "nimbus/internal/workflow/adapters/storage"
	"nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/types"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTaskService_CreateTask(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStorageInMemory()
	taskService := NewTaskService(taskStore)
	payload := "Test Payload"

	// Act
	task, err := taskService.CreateTask(payload)

	// Assert
	assert.NoError(t, err)
	assert.IsType(t, uuid.UUID{}, task.ID)
	assert.Equal(t, payload, task.Payload)
	assert.Equal(t, entity.StatusNew, task.Status)
}

func TestTaskService_GetTasks(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStorageInMemory()
	taskService := NewTaskService(taskStore)
	task := &entity.Task{
		ID:      uuid.New(),
		Payload: "Test Payload",
		Status:  entity.StatusNew,
	}

	taskStore.StoreTask(task)

	// Act
	tasks := taskService.GetTasks()

	// Assert
	assert.Len(t, tasks, 1)
	assert.Equal(t, task.Payload, tasks[0].Payload)
	assert.Equal(t, task.Status, tasks[0].Status)
}

func TestTaskService_GetTask(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStorageInMemory()
	taskService := NewTaskService(taskStore)
	task := &entity.Task{
		ID:      uuid.New(),
		Payload: "Test Payload",
		Status:  entity.StatusNew,
	}

	taskStore.StoreTask(task)

	// Act
	retrievedTask, err := taskService.GetTask(task.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, task.Payload, retrievedTask.Payload)
	assert.Equal(t, task.Status, retrievedTask.Status)
}

func TestTaskService_GetTask_NotFound(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStorageInMemory()
	taskService := NewTaskService(taskStore)
	nonExistentID := uuid.New()

	// Act
	retrievedTask, err := taskService.GetTask(nonExistentID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrievedTask)
}

func TestTaskService_StartTask(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStorageInMemory()
	taskService := NewTaskService(taskStore)
	existingTask := &entity.Task{
		ID:      uuid.New(),
		Payload: "Test Payload",
		Status:  entity.StatusNew,
	}
	taskStore.StoreTask(existingTask)

	// Act
	err := taskService.StartTask(existingTask.ID)

	// Assert
	assert.Nil(t, err)
}

func TestTaskService_StartTask_InvalidTaskStatus(t *testing.T) {
	// Arrange
	tests := []entity.Task{
		{
			ID:      uuid.New(),
			Payload: "test status in_progress",
			Status:  entity.StatusInProgress,
		},
		{
			ID:      uuid.New(),
			Payload: "test status failed",
			Status:  entity.StatusFailed,
		},
		{
			ID:      uuid.New(),
			Payload: "test status completed",
			Status:  entity.StatusCompleted,
		},
	}

	for _, existingTask := range tests {
		t.Run(existingTask.Payload, func(t *testing.T) {
			// Arrange
			taskStore := store.NewTaskStorageInMemory()
			taskService := NewTaskService(taskStore)
			taskStore.StoreTask(&existingTask)

			// Act
			err := taskService.StartTask(existingTask.ID)

			// Assert
			assert.IsType(t, &types.UnprocessableEntityError{}, err)
		})
	}
}

func TestTaskService_StartTask_NotFound(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStorageInMemory()
	taskService := NewTaskService(taskStore)

	// Act
	err := taskService.StartTask(uuid.New())

	// Assert
	assert.IsType(t, &types.RecordNotFoundError{}, err)
}

func TestTaskService_CompleteTask(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStorageInMemory()
	taskService := NewTaskService(taskStore)
	inProgressTask := &entity.Task{
		ID:      uuid.New(),
		Payload: "Test Payload",
		Status:  entity.StatusInProgress,
	}
	taskStore.StoreTask(inProgressTask)

	// Act
	err := taskService.CompleteTask(inProgressTask.ID, " - Additional Payload")
	task, _ := taskService.GetTask(inProgressTask.ID)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, entity.StatusCompleted, task.Status)
	assert.Equal(t, "Test Payload - Additional Payload", task.Payload)
}

func TestTaskService_CompleteTask_InvalidTaskStatus(t *testing.T) {
	// Arrange
	tests := []entity.Task{
		{
			ID:      uuid.New(),
			Payload: "test status new",
			Status:  entity.StatusNew,
		},
		{
			ID:      uuid.New(),
			Payload: "test status failed",
			Status:  entity.StatusFailed,
		},
		{
			ID:      uuid.New(),
			Payload: "test status completed",
			Status:  entity.StatusCompleted,
		},
	}

	for _, existingTask := range tests {
		t.Run(existingTask.Payload, func(t *testing.T) {
			// Arrange
			taskStore := store.NewTaskStorageInMemory()
			taskService := NewTaskService(taskStore)
			taskStore.StoreTask(&existingTask)

			// Act
			err := taskService.CompleteTask(existingTask.ID, "")

			// Assert
			assert.IsType(t, &types.UnprocessableEntityError{}, err)
		})
	}
}

func TestTaskService_CompleteTask_NotFound(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStorageInMemory()
	taskService := NewTaskService(taskStore)

	// Act
	err := taskService.CompleteTask(uuid.New(), "")

	// Assert
	assert.IsType(t, &types.RecordNotFoundError{}, err)
}

func TestTaskService_FailTask(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStorageInMemory()
	taskService := NewTaskService(taskStore)
	inProgressTask := &entity.Task{
		ID:      uuid.New(),
		Payload: "Test Payload",
		Status:  entity.StatusInProgress,
	}
	taskStore.StoreTask(inProgressTask)

	// Act
	err := taskService.FailTask(inProgressTask.ID, "Some failure reason")
	task, _ := taskService.GetTask(inProgressTask.ID)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, entity.StatusFailed, task.Status)
	assert.Equal(t, "Some failure reason", task.FailReason)
}

func TestTaskService_FailTask_NotFound(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStorageInMemory()
	taskService := NewTaskService(taskStore)

	// Act
	err := taskService.FailTask(uuid.New(), "Some failure reason")

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
			taskStore := store.NewTaskStorageInMemory()
			taskService := NewTaskService(taskStore)
			taskStore.StoreTask(&existingTask)

			// Act
			err := taskService.FailTask(existingTask.ID, "Some failure reason")

			// Assert
			assert.IsType(t, &types.UnprocessableEntityError{}, err)
		})
	}
}
