package sqs

import (
	"testing"

	"nimbus/internal/task_runnner/domain/entity"

	"github.com/stretchr/testify/assert"
)

func validConfig() entity.TaskRunnerConfig {
	return entity.TaskRunnerConfig{
		"queue_url":  "http://localhost:9324/000000000000/nimbus-tasks",
		"region":     "us-east-1",
		"access_key": "test",
		"secret_key": "test",
	}
}

func TestConfigValidator_ValidConfig(t *testing.T) {
	v := NewConfigValidator()
	err := v.Validate(validConfig())
	assert.NoError(t, err)
}

func TestConfigValidator_ValidConfigWithEndpoint(t *testing.T) {
	v := NewConfigValidator()
	cfg := validConfig()
	cfg["endpoint"] = "http://localhost:9324"
	err := v.Validate(cfg)
	assert.NoError(t, err)
}

func TestConfigValidator_MissingRequiredFields(t *testing.T) {
	v := NewConfigValidator()
	err := v.Validate(entity.TaskRunnerConfig{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "queue_url")
	assert.Contains(t, err.Error(), "region")
	assert.Contains(t, err.Error(), "access_key")
	assert.Contains(t, err.Error(), "secret_key")
}

func TestConfigValidator_EmptyRequiredFields(t *testing.T) {
	v := NewConfigValidator()
	cfg := entity.TaskRunnerConfig{
		"queue_url":  "  ",
		"region":     "",
		"access_key": "",
		"secret_key": "",
	}
	err := v.Validate(cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "queue_url")
	assert.Contains(t, err.Error(), "region")
	assert.Contains(t, err.Error(), "access_key")
	assert.Contains(t, err.Error(), "secret_key")
}

func TestConfigValidator_InvalidQueueURL(t *testing.T) {
	v := NewConfigValidator()
	cfg := validConfig()
	cfg["queue_url"] = "not-a-url"
	err := v.Validate(cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid queue_url")
}

func TestConfigValidator_InvalidEndpoint(t *testing.T) {
	v := NewConfigValidator()
	cfg := validConfig()
	cfg["endpoint"] = "not-a-url"
	err := v.Validate(cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid endpoint")
}

func TestConfigValidator_EmptyConfig(t *testing.T) {
	v := NewConfigValidator()
	err := v.Validate(nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required config fields")
}
