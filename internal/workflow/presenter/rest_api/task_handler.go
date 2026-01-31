package restapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"nimbus/internal/workflow/domain/entity"
	"nimbus/internal/workflow/domain/service"
	"nimbus/internal/workflow/domain/types"
)

type taskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *taskHandler {
	return &taskHandler{
		service: service,
	}
}

func (t *taskHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/tasks", t.handleCreateTask)
	rg.GET("/tasks", t.handGetTasks)
	rg.GET("/tasks/:id", t.handleGetTask)
	rg.POST("/tasks/:id/start", t.handleStartTask)
	rg.POST("/tasks/:id/complete", t.handleCompleteTask)
	rg.POST("/tasks/:id/fail", t.handleFailTask)
}

func (t *taskHandler) handleCreateTask(ctx *gin.Context) {
	var task entity.Task
	err := ctx.ShouldBindJSON(&task)
	if err != nil {
		log.Printf("%s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	taskEntity, err := t.service.CreateTask(&task)
	if err != nil {
		log.Printf("Error storing task: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	fmt.Printf("Created Task: %v", taskEntity)
	ctx.JSON(http.StatusCreated, gin.H{"data": taskEntity})
}

func (t *taskHandler) handGetTasks(ctx *gin.Context) {
	tasks := t.service.GetTasks()
	ctx.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (t *taskHandler) handleGetTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("Invalid task ID: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid task ID"})
		return
	}

	task, err := t.service.GetTask(id)
	if err != nil {
		log.Printf("Error getting task: %s", err)
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": task})
}

func (t *taskHandler) handleStartTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("Invalid task ID: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid task ID"})
		return
	}

	err = t.service.StartTask(id)
	if err != nil {
		log.Printf("Error starting task: %s", err)
		switch err.(type) {
		case *types.RecordNotFoundError:
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		case *types.UnprocessableEntityError:
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task started"})
}

func (t *taskHandler) handleCompleteTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("Invalid task ID: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid task ID"})
		return
	}

	// mock data for completion
	completionPayload := ""
	err = t.service.CompleteTask(id, completionPayload)
	if err != nil {
		log.Printf("Error completing task: %s", err)
		switch err.(type) {
		case *types.RecordNotFoundError:
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		case *types.UnprocessableEntityError:
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task completed"})
}

func (t *taskHandler) handleFailTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("Invalid task ID: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid task ID"})
		return
	}

	// mock reason for failure
	reason := "Simulated failure"
	err = t.service.FailTask(id, reason)
	if err != nil {
		log.Printf("Error failing task: %s", err)
		switch err.(type) {
		case *types.RecordNotFoundError:
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		case *types.UnprocessableEntityError:
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task marked as failed"})
}
