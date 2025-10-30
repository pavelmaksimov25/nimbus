package main

import (
	"nimbus/internal/workflow/adapters/storage"
	"nimbus/internal/workflow/application/task"
	"nimbus/internal/workflow/application/workflow"
	restapi "nimbus/internal/workflow/presenter/rest_api"
)

func main() {
	taskStore := storage.NewTaskStorageInMemory()
	taskService := task.NewTaskService(taskStore)

	workflowStore := storage.NewWorkflowInMemoryStorage()
	workflowService := workflow.NewWorkflowService(workflowStore)

	server := restapi.NewRestApiServer()
	server.RegisterRoutes(taskService, workflowService)
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
