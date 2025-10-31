package main

import (
	"nimbus/internal/workflow/adapters/storage"
	"nimbus/internal/workflow/application/task"

	workflow_repository "nimbus/internal/workflow/adapters/repository/workflow"
	workflow_service "nimbus/internal/workflow/application/workflow"
	restapi "nimbus/internal/workflow/presenter/rest_api"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	taskStore := storage.NewTaskStorageInMemory()
	taskService := task.NewTaskService(taskStore)

	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local" // todo :: move to env

	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	workflowRepository := workflow_repository.NewWorkflowRepository(dbConn)
	workflowService := workflow_service.NewWorkflowService(workflowRepository)

	server := restapi.NewRestApiServer()
	server.RegisterRoutes(taskService, workflowService)
	server.Run()
	if err != nil {
		panic(err)
	}
}
