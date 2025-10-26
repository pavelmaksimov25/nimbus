package entity

import (
	"time"

	uuid "github.com/google/uuid"
)

type Workflow struct {
	ID      uuid.UUID
	Name    string
	Tasks   []Task
	CreatedAt time.Time
}

type Task struct {
	ID      uuid.UUID
	Payload string
	CreatedAt time.Time
}
