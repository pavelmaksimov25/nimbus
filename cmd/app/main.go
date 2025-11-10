package main

import (
	"log"
	"nimbus/internal/workflow/adapters/storage"
	"nimbus/internal/workflow/application/task"
	"os"

	workflow_repository "nimbus/internal/workflow/adapters/repository/workflow"
	workflow_service "nimbus/internal/workflow/application/workflow"
	restapi "nimbus/internal/workflow/presenter/rest_api"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	taskStore := storage.NewTaskStorageInMemory()
	taskService := task.NewTaskService(taskStore)

	dsn := os.Getenv("DB_DSN")

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
