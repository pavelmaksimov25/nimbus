package workflow

import (
	"errors"
	"testing"

	"nimbus/internal/workflow/adapters/mocks"
	"nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/types"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWorkflowService_CreateWorkflow(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWorkflowRepository(ctrl)
	workflowService := NewWorkflowService(mockRepo)

	mockRepo.EXPECT().
		StoreWorkflow(gomock.Any()).
		DoAndReturn(func(workflow *entity.Workflow) (*entity.Workflow, error) {
			return workflow, nil
		})

	// Act
	result, err := workflowService.CreateWorkflow("Test Workflow")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Workflow", result.Name)
	assert.IsType(t, uuid.UUID{}, result.ID)
}

func TestWorkflowService_CreateWorkflow_RepositoryError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWorkflowRepository(ctrl)
	workflowService := NewWorkflowService(mockRepo)

	mockRepo.EXPECT().
		StoreWorkflow(gomock.Any()).
		Return(nil, errors.New("database error"))

	// Act
	result, err := workflowService.CreateWorkflow("Test Workflow")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestWorkflowService_GetWorkflows(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWorkflowRepository(ctrl)
	workflowService := NewWorkflowService(mockRepo)

	expectedWorkflow := entity.Workflow{
		ID:   uuid.New(),
		Name: "Test Workflow",
	}

	mockRepo.EXPECT().
		GetWorkflows().
		Return([]entity.Workflow{expectedWorkflow})

	// Act
	workflows := workflowService.GetWorkflows()

	// Assert
	assert.Len(t, workflows, 1)
	assert.Equal(t, expectedWorkflow.Name, workflows[0].Name)
	assert.Equal(t, expectedWorkflow.ID, workflows[0].ID)
}

func TestWorkflowService_GetWorkflows_Empty(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWorkflowRepository(ctrl)
	workflowService := NewWorkflowService(mockRepo)

	mockRepo.EXPECT().
		GetWorkflows().
		Return([]entity.Workflow{})

	// Act
	workflows := workflowService.GetWorkflows()

	// Assert
	assert.Empty(t, workflows)
}

func TestWorkflowService_GetWorkflow(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWorkflowRepository(ctrl)
	workflowService := NewWorkflowService(mockRepo)

	expectedWorkflow := &entity.Workflow{
		ID:   uuid.New(),
		Name: "Test Workflow",
	}

	mockRepo.EXPECT().
		GetWorkflow(expectedWorkflow.ID).
		Return(expectedWorkflow)

	// Act
	result, err := workflowService.GetWorkflow(expectedWorkflow.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedWorkflow.Name, result.Name)
	assert.Equal(t, expectedWorkflow.ID, result.ID)
}

func TestWorkflowService_GetWorkflow_NotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockWorkflowRepository(ctrl)
	workflowService := NewWorkflowService(mockRepo)
	nonExistentID := uuid.New()

	mockRepo.EXPECT().
		GetWorkflow(nonExistentID).
		Return(nil)

	// Act
	result, err := workflowService.GetWorkflow(nonExistentID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &types.RecordNotFoundError{}, err)
}
