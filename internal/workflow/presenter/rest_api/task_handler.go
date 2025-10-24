package restapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	task "nimbus/internal/workflow/application/task"
)

type Task struct {
	ID    string     `json:"id,omitempty"`
	Payload string 	 `json:"payload"`
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
	rg.GET("/tasks/:id/status", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Nope"})
	})
	rg.GET("/tasks/:id", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Nope"})
	})
}

func (t *taskHandler) handleCreateTask(ctx *gin.Context) {
	var task Task
	err := ctx.ShouldBindJSON(&task)
	if err != nil {
		log.Printf("%s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}



	taskID, err := t.service.CreateTask(task.Payload)
	if err != nil {
		log.Printf("Error storing task: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	task.ID = taskID

	fmt.Printf("Received Task: %v", task)
	ctx.JSON(http.StatusCreated, gin.H{"data": task})
}

func (t *taskHandler) handGetTasks(ctx *gin.Context) {
	tasks := t.service.GetTasks()
	ctx.JSON(http.StatusOK, gin.H{"data": tasks})
}