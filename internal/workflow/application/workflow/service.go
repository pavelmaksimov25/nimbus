package workflow

import (
	"nimbus/internal/workflow/adapters/storage"
	"nimbus/internal/workflow/domain/entity"
	"time"

	"github.com/google/uuid"
)


type WorkflowService interface {
	CreateWorkflow(name string) (*entity.Workflow, error)
}

type workflowService struct {
	storage storage.WorkflowStorage
}

func NewWorkflowService(storage storage.WorkflowStorage) WorkflowService {
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
