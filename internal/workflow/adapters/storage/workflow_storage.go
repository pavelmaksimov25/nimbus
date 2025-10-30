package storage

import (
	"nimbus/internal/workflow/domain/entity"
	"sync"

	"github.com/google/uuid"
)

type workflowStorage struct {
	mu    sync.RWMutex
	store map[uuid.UUID]*entity.Workflow
}

func NewWorkflowInMemoryStorage() WorkflowStorage {
	return &workflowStorage{
		store: make(map[uuid.UUID]*entity.Workflow),
	}
}

func (ws *workflowStorage) StoreWorkflow(workflow *entity.Workflow) (*entity.Workflow, error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.store[workflow.ID] = workflow
	return workflow, nil
}

func (ws *workflowStorage) GetWorkflows() []entity.Workflow {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	var workflows []entity.Workflow
	for _, workflow := range ws.store {
		workflows = append(workflows, *workflow)
	}
	return workflows
}

func (ws *workflowStorage) GetWorkflow(id uuid.UUID) *entity.Workflow {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	workflow, exists := ws.store[id]
	if !exists {
		return nil
	}
	return workflow
}

func (ws *workflowStorage) UpdateWorkflow(workflow *entity.Workflow) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.store[workflow.ID] = workflow
	return nil
}
