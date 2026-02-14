package restapi

import (
	"log"
	"net/http"

	"nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/service"
	"nimbus/internal/workflow/domain/types"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type workflowHandler struct {
	service service.WorkflowService
}

func NewWorkflowHandler(service service.WorkflowService) *workflowHandler {
	return &workflowHandler{
		service: service,
	}
}

func (w *workflowHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/workflows", w.handleCreateWorkflow)
	rg.GET("/workflows", w.handleGetWorkflows)
	rg.GET("/workflows/:id", w.handleGetWorkflow)
}

func (w *workflowHandler) handleCreateWorkflow(ctx *gin.Context) {
	var workflow entity.Workflow
	if err := ctx.ShouldBindJSON(&workflow); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	storedWorkflow, err := w.service.CreateWorkflow(workflow.Name)
	if err != nil {
		log.Printf("Error creating workflow: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": storedWorkflow})
}

func (w *workflowHandler) handleGetWorkflows(ctx *gin.Context) {
	workflows := w.service.GetWorkflows()
	ctx.JSON(http.StatusOK, gin.H{"data": workflows})
}

func (w *workflowHandler) handleGetWorkflow(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("Invalid workflow ID: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid workflow ID"})
		return
	}

	workflow, err := w.service.GetWorkflow(id)
	if err != nil {
		log.Printf("Error getting workflow: %s", err)
		switch err.(type) {
		case *types.RecordNotFoundError:
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"data": workflow})
}
