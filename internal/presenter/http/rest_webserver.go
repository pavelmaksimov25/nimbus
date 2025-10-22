package restserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const routePrefix string = "/api/v1"

type Task struct {
	ID string `json:"id,omitempty"`
	Items []TaskItem `json:"items"`
}

type TaskItem struct {
	ID string `json:"id"`
}

type TaskWebServer struct {
	router *gin.Engine
	addr string
}

func NewTaskWebServer() *TaskWebServer {
	return &TaskWebServer{
		router: gin.Default(),
		addr: ":8080",
	}
}

func (t *TaskWebServer) RegisterRoutes() {
	rg := t.router.Group(routePrefix)
	rg.POST("/tasks", t.handleCreateTask)
	rg.GET("/tasks/:id/status", func (c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Nope"})
	})
	rg.GET("/tasks/:id", func (c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Nope"})
	})
}

func (t *TaskWebServer) Run() error {
	return t.router.Run(t.addr)
}

func (t *TaskWebServer) handleCreateTask(ctx *gin.Context) {
	var task Task
	err := ctx.ShouldBindJSON(&task); if err != nil {
		log.Printf("%s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	task.ID = "test"

	fmt.Printf("Received Task: %v", task)
	ctx.JSON(http.StatusCreated, gin.H{"data": task})
}

