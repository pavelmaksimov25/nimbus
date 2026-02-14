package echo

import (
	"context"
	"testing"

	"nimbus/internal/task_runnner/domain/entity"

	"github.com/stretchr/testify/assert"
)

func TestEchoRunner_Execute(t *testing.T) {
	factory := NewFactory()
	r := factory(entity.TaskRunnerConfig{})

	err := r.Execute(context.Background(), `{"action": "test"}`)

	assert.NoError(t, err)
}
