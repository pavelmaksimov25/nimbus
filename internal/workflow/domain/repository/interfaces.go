package repository

type TaskRepository interface {
	Save(taskID string, payload string) error
	FindByID(taskID string) (string, error)
}
