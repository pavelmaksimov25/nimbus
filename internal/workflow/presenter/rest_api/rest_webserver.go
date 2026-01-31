package restapi

import (
	"net/http"

	trService "nimbus/internal/task_runnner/domain/service"
	trRestApi "nimbus/internal/task_runnner/presenter/rest_api"
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

func maxBodySize() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20) // 1 MB
		c.Next()
	}
}

func (t *ApiServer) RegisterRoutes(taskService service.TaskService, workflowService service.WorkflowService, taskRunnerService trService.TaskRunnerService) {
	t.router.Use(maxBodySize())
	rg := t.router.Group(routePrefix)

	taskHandler := NewTaskHandler(taskService)
	taskHandler.RegisterRoutes(rg)

	workflowHandler := NewWorkflowHandler(workflowService)
	workflowHandler.RegisterRoutes(rg)

	taskRunnerHandler := trRestApi.NewTaskRunnerHandler(taskRunnerService)
	taskRunnerHandler.RegisterRoutes(rg)
}

func (t *ApiServer) Run() error {
	return t.router.Run(t.addr)
}
