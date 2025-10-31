package restapi

import (
	"nimbus/internal/workflow/domain/service"

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

func (t *ApiServer) RegisterRoutes(taskService service.TaskService, workflowService service.WorkflowService) {
	rg := t.router.Group(routePrefix)

	taskHandler := NewTaskHandler(taskService)
	taskHandler.RegisterRoutes(rg)

	workflowHandler := NewWorkflowHandler(workflowService)
	workflowHandler.RegisterRoutes(rg)
}

func (t *ApiServer) Run() error {
	return t.router.Run(t.addr)
}
