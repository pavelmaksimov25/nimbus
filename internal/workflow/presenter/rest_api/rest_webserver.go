package restapi

import (
	"nimbus/internal/workflow/application/task"

	"github.com/gin-gonic/gin"
)

const routePrefix string = "/api/v1"

type ApiServer struct {
	router *gin.Engine
	addr   string
}

func NewRestApiServer() *ApiServer {
	return &ApiServer{
		router: gin.Default(),
		addr:   ":8080",
	}
}

func (t *ApiServer) RegisterRoutes(taskService task.Service) {
	rg := t.router.Group(routePrefix)

	taskHandler := NewTaskHandler(taskService)
	taskHandler.RegisterRoutes(rg)
}

func (t *ApiServer) Run() error {
	return t.router.Run(t.addr)
}
