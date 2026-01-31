package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaskRunnerType string

func (trt *TaskRunnerType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*trt = TaskRunnerType(v)
	case string:
		*trt = TaskRunnerType(v)
	default:
		return fmt.Errorf("unsupported type for TaskRunnerType: %T", value)
	}
	return nil
}

func (trt *TaskRunnerType) Value() (driver.Value, error) {
	return string(*trt), nil
}

const (
	Queue TaskRunnerType = "queue"
)

type TaskRunnerConfig map[string]interface{}

func (c TaskRunnerConfig) Value() (driver.Value, error) {
	valueString, err := json.Marshal(c)
	return string(valueString), err
}

func (c *TaskRunnerConfig) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported type for TaskRunnerConfig: %T", value)
	}
	return json.Unmarshal(bytes, c)
}

type TaskRunner struct {
	ID        uuid.UUID        `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string           `json:"name" gorm:"type:varchar(255);not null"`
	Type      TaskRunnerType   `json:"type" gorm:"type:task_runner_type;default:'queue'"`
	Config    TaskRunnerConfig `json:"config" gorm:"type:jsonb;default:'{}'"`
	CreatedAt time.Time        `json:"created_at" gorm:"autoCreateTime"`
}

func (TaskRunner) TableName() string {
	return "task_runners"
}

type TaskRunnerTask struct {
	TaskRunnerID uuid.UUID `json:"task_runner_id" gorm:"type:uuid;primaryKey"`
	TaskID       uuid.UUID `json:"task_id" gorm:"type:uuid;primaryKey"`
}

func (TaskRunnerTask) TableName() string {
	return "task_runners_tasks"
}
