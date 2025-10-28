package entity

import (
	"time"

	uuid "github.com/google/uuid"
)

type Workflow struct {
	ID        uuid.UUID
	Name      string
	Tasks     []Task
	CreatedAt time.Time
}

type Task struct {
	ID         uuid.UUID  `json:"id,omitempty"`
	Payload    string     `json:"payload"`
	Status     TaskStatus `json:"status,omitempty"`
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
