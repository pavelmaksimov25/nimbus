package task

import (
	entity "nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type taskRepository struct {
	conn *gorm.DB
}

func NewTaskRepository(conn *gorm.DB) repository.TaskRepository {
	return &taskRepository{
		conn: conn,
	}
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
