package workflow

import (
	"nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/repository"
	"nimbus/internal/workflow/domain/service"
	"nimbus/internal/workflow/domain/types"
	"time"

	"github.com/google/uuid"
)

type workflowService struct {
	repository repository.WorkflowRepository
}

func NewWorkflowService(repository repository.WorkflowRepository) service.WorkflowService {
	return &workflowService{
		repository: repository,
	}
}

func (w *workflowService) CreateWorkflow(name string) (*entity.Workflow, error) {
	return w.repository.StoreWorkflow(&entity.Workflow{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
	})
}

func (w *workflowService) GetWorkflows() []entity.Workflow {
	return w.repository.GetWorkflows()
}

func (w *workflowService) GetWorkflow(id uuid.UUID) (*entity.Workflow, error) {
	workflow := w.repository.GetWorkflow(id)
	if workflow == nil {
		return nil, &types.RecordNotFoundError{Resource: "Workflow", ID: id.String()}
	}
	return workflow, nil
}
