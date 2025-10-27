package restapi

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	task "nimbus/internal/workflow/application/task"
)

type Task struct {
	ID        string `json:"id,omitempty"`
	Payload   string `json:"payload"`
	Status    string `json:"status,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type taskHandler struct {
	// for testing purposes the store will be used directly
	service task.Service
}

func NewTaskHandler(service task.Service) *taskHandler {
	return &taskHandler{
		service: service,
	}
}

func (t *taskHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/tasks", t.handleCreateTask)
	rg.GET("/tasks", t.handGetTasks)
	rg.GET("/tasks/:id", t.handleGetTask)
	rg.POST("/tasks/:id/start", t.handleStartTask)
}

func (t *taskHandler) handleCreateTask(ctx *gin.Context) {
	var task Task
	err := ctx.ShouldBindJSON(&task)
	if err != nil {
		log.Printf("%s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	taskEntity, err := t.service.CreateTask(task.Payload)
	if err != nil {
		log.Printf("Error storing task: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	task.ID = taskEntity.ID.String()
	task.CreatedAt = taskEntity.CreatedAt.Format(time.RFC3339)
	task.Status = string(taskEntity.Status)

	fmt.Printf("Received Task: %v", task)
	ctx.JSON(http.StatusCreated, gin.H{"data": task})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task started"})
}