package workflow

import (
	"nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type workflowRepository struct {
	conn *gorm.DB
}

func NewWorkflowRepository(conn *gorm.DB) repository.WorkflowRepository {
	return &workflowRepository{
		conn: conn,
	}
}

func (w *workflowRepository) StoreWorkflow(workflow *entity.Workflow) (*entity.Workflow, error) {
	err := w.conn.Create(&workflow).Error
	return workflow, err
}

func (w *workflowRepository) GetWorkflows() []entity.Workflow {
	// Implementation goes here
	return nil
}

func (w *workflowRepository) GetWorkflow(id uuid.UUID) *entity.Workflow {
	// Implementation goes here
	return nil
}

func (w *workflowRepository) UpdateWorkflow(workflow *entity.Workflow) error {
	// Implementation goes here
	return nil
}
