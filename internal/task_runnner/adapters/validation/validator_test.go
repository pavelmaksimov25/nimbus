package validation

import (
	"testing"

	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/runner"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidator_AllRequiredPresent(t *testing.T) {
	v := NewConfigValidator([]runner.FieldRule{
		{Name: "a", Required: true},
		{Name: "b", Required: true},
	})
	err := v.Validate(entity.TaskRunnerConfig{"a": "x", "b": "y"})
	assert.NoError(t, err)
}

func TestConfigValidator_MissingRequired(t *testing.T) {
	v := NewConfigValidator([]runner.FieldRule{
		{Name: "a", Required: true},
		{Name: "b", Required: true},
	})
	err := v.Validate(entity.TaskRunnerConfig{"a": "x"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "b")
}

func TestConfigValidator_EmptyRequired(t *testing.T) {
	v := NewConfigValidator([]runner.FieldRule{
		{Name: "a", Required: true},
	})
	err := v.Validate(entity.TaskRunnerConfig{"a": "  "})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "a")
}

func TestConfigValidator_ValidURL(t *testing.T) {
	v := NewConfigValidator([]runner.FieldRule{
		{Name: "endpoint", Required: true, URL: true},
	})
	err := v.Validate(entity.TaskRunnerConfig{"endpoint": "http://localhost:9324"})
	assert.NoError(t, err)
}

func TestConfigValidator_InvalidURL(t *testing.T) {
	v := NewConfigValidator([]runner.FieldRule{
		{Name: "endpoint", Required: true, URL: true},
	})
	err := v.Validate(entity.TaskRunnerConfig{"endpoint": "not-a-url"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid endpoint")
}

func TestConfigValidator_OptionalURLSkippedWhenEmpty(t *testing.T) {
	v := NewConfigValidator([]runner.FieldRule{
		{Name: "endpoint", Required: false, URL: true},
	})
	err := v.Validate(entity.TaskRunnerConfig{})
	assert.NoError(t, err)
}

func TestConfigValidator_OptionalURLValidatedWhenPresent(t *testing.T) {
	v := NewConfigValidator([]runner.FieldRule{
		{Name: "endpoint", Required: false, URL: true},
	})
	err := v.Validate(entity.TaskRunnerConfig{"endpoint": "not-a-url"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid endpoint")
}

func TestConfigValidator_NoRules(t *testing.T) {
	v := NewConfigValidator([]runner.FieldRule{})
	err := v.Validate(entity.TaskRunnerConfig{"anything": "goes"})
	assert.NoError(t, err)
}

func TestConfigValidator_NilConfig(t *testing.T) {
	v := NewConfigValidator([]runner.FieldRule{
		{Name: "a", Required: true},
	})
	err := v.Validate(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required config fields")
}
