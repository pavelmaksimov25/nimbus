package task

import (
	"testing"

	store "nimbus/internal/workflow/adapters/store"

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
}

func TestTaskService_GetTasks(t *testing.T) {
	// Arrange
	taskStore := store.NewTaskStoreInMemory()
	taskService := NewTaskService(taskStore)
	payload := "Test Payload"

	taskStore.StoreTask(payload)

	// Act
	tasks := taskService.GetTasks()

	// Assert
	assert.Len(t, tasks, 1)
	assert.Equal(t, payload, tasks[0].Payload)
}
