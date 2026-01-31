package taskrunner

import (
	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type taskRunnerRepository struct {
	conn *gorm.DB
}

func NewTaskRunnerRepository(conn *gorm.DB) repository.TaskRunnerRepository {
	return &taskRunnerRepository{conn: conn}
}

func (r *taskRunnerRepository) Store(runner *entity.TaskRunner) (*entity.TaskRunner, error) {
	err := r.conn.Create(runner).Error
	return runner, err
}

func (r *taskRunnerRepository) GetAll() ([]entity.TaskRunner, error) {
	var runners []entity.TaskRunner
	err := r.conn.Find(&runners).Error
	return runners, err
}

func (r *taskRunnerRepository) GetByID(id uuid.UUID) (*entity.TaskRunner, error) {
	var runner entity.TaskRunner
	err := r.conn.First(&runner, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &runner, nil
}

func (r *taskRunnerRepository) AssignTask(runnerID uuid.UUID, taskID uuid.UUID) error {
	assignment := entity.TaskRunnerTask{
		TaskRunnerID: runnerID,
		TaskID:       taskID,
	}
	return r.conn.Create(&assignment).Error
}

func (r *taskRunnerRepository) UnassignTask(runnerID uuid.UUID, taskID uuid.UUID) error {
	return r.conn.Where("task_runner_id = ? AND task_id = ?", runnerID, taskID).
		Delete(&entity.TaskRunnerTask{}).Error
}

func (r *taskRunnerRepository) GetByTaskID(taskID uuid.UUID) ([]entity.TaskRunner, error) {
	var runners []entity.TaskRunner
	err := r.conn.
		Joins("JOIN task_runners_tasks ON task_runners_tasks.task_runner_id = task_runners.id").
		Where("task_runners_tasks.task_id = ?", taskID).
		Find(&runners).Error
	return runners, err
}
