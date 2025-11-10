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
	var workflows []entity.Workflow
	w.conn.Find(&workflows)
	return workflows
}

func (w *workflowRepository) GetWorkflow(id uuid.UUID) *entity.Workflow {
	var workflow entity.Workflow
	result := w.conn.First(&workflow, "id = ?", id)
	if result.Error != nil {
		return nil
	}
	return &workflow
}

func (w *workflowRepository) UpdateWorkflow(workflow *entity.Workflow) error {
	return w.conn.Save(workflow).Error
}
