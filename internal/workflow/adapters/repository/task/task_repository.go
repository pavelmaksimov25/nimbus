package task

import (
	"fmt"

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
	err := t.conn.Create(&task).Error
	return task, err
}

func (t *taskRepository) GetTasks() []entity.Task {
	var tasks []entity.Task
	t.conn.Find(&tasks)
	return tasks
}

func (t *taskRepository) GetTask(id uuid.UUID) *entity.Task {
	var task entity.Task
	result := t.conn.First(&task, "id = ?", id)
	if result.Error != nil {
		return nil
	}
	return &task
}

func (t *taskRepository) UpdateTask(task *entity.Task) error {
	return t.conn.Save(task).Error
}

func (t *taskRepository) UpdateTaskStatus(id uuid.UUID, fromStatus, toStatus entity.TaskStatus) error {
	result := t.conn.Model(&entity.Task{}).
		Where("id = ? AND status = ?", id, fromStatus).
		Update("status", toStatus)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task %s is not in status %s", id, fromStatus)
	}
	return nil
}
