package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONB: %T", value)
	}
	return json.Unmarshal(bytes, j)
}

type Workflow struct {
	ID        uuid.UUID `json:"id,omitempty" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
}

type Task struct {
	ID         uuid.UUID  `json:"id,omitempty" gorm:"type:uuid;default:uuid_generate_v4()"`
	Payload    string     `json:"payload" gorm:"type:jsonb"`
	WorkflowID uuid.UUID  `json:"workflow_id" gorm:"type:uuid"`
	Status     TaskStatus `json:"status" gorm:"type:TaskStatus;default:'NEW'"`
	CreatedAt  time.Time  `json:"created_at,omitempty" gorm:"autoCreateTime"`
	FailReason string     `json:"fail_reason,omitempty" gorm:"type:text"`
}

type TaskStatus string

const (
	StatusNew        TaskStatus = "NEW"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusCompleted  TaskStatus = "COMPLETED"
	StatusFailed     TaskStatus = "FAILED"
)

func (ts *TaskStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*ts = TaskStatus(v)
	case string:
		*ts = TaskStatus(v)
	default:
		return fmt.Errorf("unsupported type for TaskStatus: %T", value)
	}
	return nil
}

func (ts *TaskStatus) Value() (driver.Value, error) {
	return string(*ts), nil
}
