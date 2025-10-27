package types

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("record not found")
var ErrUnpocessable = errors.New("record unprocessable")

type RecordNotFoundError struct {
	Resource string
	ID       string
	Msg      string
}

func (e *RecordNotFoundError) Error() string {
	baseMsg := fmt.Sprintf("%s with ID '%s' was not found", e.Resource, e.ID)
	if e.Msg != "" {
		return fmt.Sprintf("%s: %s", baseMsg, e.Msg)
	}
	return baseMsg
}

func (e *RecordNotFoundError) Unwrap() error {
	return ErrNotFound
}

type UnprocessableEntityError struct {
	Resource string
	ID       string
	Msg      string
}

func (e *UnprocessableEntityError) Error() string {
	baseMsg := fmt.Sprintf("%s with ID '%s' is unprocessable", e.Resource, e.ID)
	if e.Msg != "" {
		return fmt.Sprintf("%s: %s", baseMsg, e.Msg)
	}
	return baseMsg
}

func (e *UnprocessableEntityError) Unwrap() error {
	return ErrUnpocessable
}
