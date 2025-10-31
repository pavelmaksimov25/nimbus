package task

import (
	entity "nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/repository"

	"github.com/google/uuid"
)

type taskRepository struct {
	store repository.TaskRepository
}

func NewTaskService() repository.TaskRepository {
	return &taskRepository{}
}

func (t *taskRepository) StoreTask(task *entity.Task) (*entity.Task, error) {
	// Implementation goes here
	return nil, nil
}

func (t *taskRepository) GetTasks() []entity.Task {
	// Implementation goes here
	return nil
}

func (t *taskRepository) GetTask(id uuid.UUID) *entity.Task {
	// Implementation goes here
	return nil
}

func (t *taskRepository) UpdateTask(task *entity.Task) error {
	// Implementation goes here
	return nil
}