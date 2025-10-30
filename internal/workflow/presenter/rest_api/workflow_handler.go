package restapi

import (
	"nimbus/internal/workflow/application/workflow"
	"nimbus/internal/workflow/domain/entity"

	"github.com/gin-gonic/gin"
)

type workflowHandler struct {
	service workflow.WorkflowService
}


func NewWorkflowHandler(service workflow.WorkflowService) *workflowHandler {
	return &workflowHandler{
		service: service,
	}
}

func (w *workflowHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/workflows", w.handleCreateWorkflow)
}

func (w *workflowHandler) handleCreateWorkflow(ctx *gin.Context) {
	var workflow entity.Workflow
	if err := ctx.ShouldBindJSON(&workflow); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	storedWorkflow, err := w.service.CreateWorkflow(workflow.Name)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create workflow"})
		return
	}

	ctx.JSON(201, gin.H{"data": storedWorkflow})
}