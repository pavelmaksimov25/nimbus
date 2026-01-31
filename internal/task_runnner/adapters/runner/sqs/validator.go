package sqs

import (
	"nimbus/internal/task_runnner/adapters/validation"
	"nimbus/internal/task_runnner/domain/runner"
)

func NewConfigValidator() runner.ConfigValidator {
	return validation.NewConfigValidator([]runner.FieldRule{
		{Name: "queue_url", Required: true, URL: true},
		{Name: "region", Required: true},
		{Name: "access_key", Required: true},
		{Name: "secret_key", Required: true},
		{Name: "endpoint", Required: false, URL: true},
	})
}
