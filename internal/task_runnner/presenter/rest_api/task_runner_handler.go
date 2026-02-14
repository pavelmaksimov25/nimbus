package restapi

import (
	"log"
	"net/http"
	"time"

	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/service"
	"nimbus/internal/workflow/domain/types"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var sensitiveConfigKeys = map[string]bool{
	"access_key": true,
	"secret_key": true,
}

type createRunnerRequest struct {
	Name   string                  `json:"name" binding:"required"`
	Type   entity.TaskRunnerType   `json:"type" binding:"required"`
	Config entity.TaskRunnerConfig `json:"config"`
}

type taskRunnerResponse struct {
	ID        uuid.UUID               `json:"id"`
	Name      string                  `json:"name"`
	Type      entity.TaskRunnerType   `json:"type"`
	Config    entity.TaskRunnerConfig `json:"config"`
	CreatedAt time.Time               `json:"created_at"`
}

func newTaskRunnerResponse(r *entity.TaskRunner) taskRunnerResponse {
	config := make(entity.TaskRunnerConfig, len(r.Config))
	for k, v := range r.Config {
		if sensitiveConfigKeys[k] {
			config[k] = "[REDACTED]"
			continue
		}

		config[k] = v
	}

	return taskRunnerResponse{
		ID:        r.ID,
		Name:      r.Name,
		Type:      r.Type,
		Config:    config,
		CreatedAt: r.CreatedAt,
	}
}

func newTaskRunnerListResponse(runners []entity.TaskRunner) []taskRunnerResponse {
	resp := make([]taskRunnerResponse, len(runners))
	for i := range runners {
		resp[i] = newTaskRunnerResponse(&runners[i])
	}
	return resp
}

type taskRunnerHandler struct {
	service service.TaskRunnerService
}

func NewTaskRunnerHandler(service service.TaskRunnerService) *taskRunnerHandler {
	return &taskRunnerHandler{service: service}
}

func (h *taskRunnerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/task-runners", h.handleCreateRunner)
	rg.GET("/task-runners", h.handleGetRunners)
	rg.GET("/task-runners/:id", h.handleGetRunner)
	rg.POST("/task-runners/:id/assign/:taskId", h.handleAssignTask)
	rg.DELETE("/task-runners/:id/assign/:taskId", h.handleUnassignTask)
	rg.GET("/tasks/:id/runners", h.handleGetRunnersByTask)
}

func (h *taskRunnerHandler) handleCreateRunner(ctx *gin.Context) {
	var req createRunnerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	runner := &entity.TaskRunner{
		Name:   req.Name,
		Type:   req.Type,
		Config: req.Config,
	}

	created, err := h.service.CreateRunner(runner)
	if err != nil {
		switch err.(type) {
		case *types.UnprocessableEntityError:
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		default:
			log.Printf("Error creating task runner: %s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": newTaskRunnerResponse(created)})
}

func (h *taskRunnerHandler) handleGetRunners(ctx *gin.Context) {
	runners, err := h.service.GetRunners()
	if err != nil {
		log.Printf("Error getting task runners: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newTaskRunnerListResponse(runners)})
}

func (h *taskRunnerHandler) handleGetRunner(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid runner ID"})
		return
	}

	runner, err := h.service.GetRunner(id)
	if err != nil {
		switch err.(type) {
		case *types.RecordNotFoundError:
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newTaskRunnerResponse(runner)})
}

func (h *taskRunnerHandler) handleAssignTask(ctx *gin.Context) {
	runnerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid runner ID"})
		return
	}

	taskID, err := uuid.Parse(ctx.Param("taskId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid task ID"})
		return
	}

	err = h.service.AssignTask(runnerID, taskID)
	if err != nil {
		switch err.(type) {
		case *types.RecordNotFoundError:
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		default:
			log.Printf("Error assigning task to runner: %s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task assigned to runner"})
}

func (h *taskRunnerHandler) handleUnassignTask(ctx *gin.Context) {
	runnerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid runner ID"})
		return
	}

	taskID, err := uuid.Parse(ctx.Param("taskId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid task ID"})
		return
	}

	err = h.service.UnassignTask(runnerID, taskID)
	if err != nil {
		log.Printf("Error unassigning task from runner: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task unassigned from runner"})
}

func (h *taskRunnerHandler) handleGetRunnersByTask(ctx *gin.Context) {
	taskID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid task ID"})
		return
	}

	runners, err := h.service.GetRunnersByTaskID(taskID)
	if err != nil {
		log.Printf("Error getting runners for task: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newTaskRunnerListResponse(runners)})
}
