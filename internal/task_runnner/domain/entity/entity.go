package entity

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type TaskRunnerType string

func (trt *TaskRunnerType) Scan(value interface{}) error {
	*trt = TaskRunnerType(value.([]byte))
	return nil
}

func (trt *TaskRunnerType) Value() (driver.Value, error) {
	return string(*trt), nil
}

const (
	Queue TaskRunnerType = "queue"
)

type TaskRunner struct {
	ID   uuid.UUID      `json:"id,omitempty" gorm:"type:uuid;default:uuid_generate_v4()"`
	Type TaskRunnerType `json:"type" gorm:"type:TaskRunnerType"`
}

