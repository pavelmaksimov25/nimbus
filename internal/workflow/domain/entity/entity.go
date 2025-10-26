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
	ID        uuid.UUID
	Payload   string
	Status    TaskStatus
	CreatedAt time.Time
}

type TaskStatus string

const (
	StatusNew        TaskStatus = "NEW"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusCompleted  TaskStatus = "COMPLETED"
	StatusFailed     TaskStatus = "FAILED"
)
