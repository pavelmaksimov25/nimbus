package entity

import (
	uuid "github.com/google/uuid"
)

type Workflow struct {
	ID      uuid.UUID
	Payload string
	Tasks   []Task
}

type Task struct {
	ID      uuid.UUID
	Payload string
}
