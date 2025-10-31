package workflow

import (
	"nimbus/internal/workflow/adapters/storage"
	"nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/service"
	"time"

	"github.com/google/uuid"
)

type workflowService struct {
	storage storage.WorkflowStorage
}

func NewWorkflowService(storage storage.WorkflowStorage) service.WorkflowService {
	return &workflowService{
		storage: storage,
	}
}

func (w *workflowService) CreateWorkflow(name string) (*entity.Workflow, error) {
	return w.storage.StoreWorkflow(&entity.Workflow{
		ID: uuid.New(),
		Name: name,
		CreatedAt: time.Now(),
	})
}
