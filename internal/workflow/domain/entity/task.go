package entity

import (
	uuid "github.com/google/uuid"
)

type Task struct {
	ID uuid.UUID
	Payload string
}