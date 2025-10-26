package task

import (
	"testing"

	store "nimbus/internal/workflow/adapters/store"
	"nimbus/internal/workflow/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTaskService_CreateTask(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStoreInMemory()
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
	taskStore := store.NewTaskStoreInMemory()
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
	taskStore := store.NewTaskStoreInMemory()
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
	taskStore := store.NewTaskStoreInMemory()
	taskService := NewTaskService(taskStore)
	nonExistentID := uuid.New()

	// Act
	retrievedTask, err := taskService.GetTask(nonExistentID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrievedTask)
}