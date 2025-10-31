package entity

import (
	"time"

	uuid "github.com/google/uuid"
)

type Workflow struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Task struct {
	ID         uuid.UUID  `json:"id,omitempty"`
	Payload    string     `json:"payload"`
	WorkflowID uuid.UUID  `json:"workflow_id"`
	Status     TaskStatus `json:"status"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	FailReason string     `json:"fail_reason,omitempty"`
}

type TaskStatus string

const (
	StatusNew        TaskStatus = "NEW"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusCompleted  TaskStatus = "COMPLETED"
	StatusFailed     TaskStatus = "FAILED"
)
