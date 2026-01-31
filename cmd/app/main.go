package main

import (
	"log"
	"os"

	trEntity "nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/runner"

	trRepo "nimbus/internal/task_runnner/adapters/repository/taskrunner"
	echoAdapter "nimbus/internal/task_runnner/adapters/runner/echo"
	sqsAdapter "nimbus/internal/task_runnner/adapters/runner/sqs"
	"nimbus/internal/task_runnner/adapters/validation"
	dispatchSvc "nimbus/internal/task_runnner/application/dispatch"
	trSvc "nimbus/internal/task_runnner/application/taskrunner"

	taskRepository "nimbus/internal/workflow/adapters/repository/task"
	workflowRepository "nimbus/internal/workflow/adapters/repository/workflow"
	"nimbus/internal/workflow/application/task"
	workflowService "nimbus/internal/workflow/application/workflow"
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

	dsn := os.Getenv("DB_DSN")
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Runner factories
	factories := map[trEntity.TaskRunnerType]runner.Factory{
		trEntity.AwsSqs: sqsAdapter.NewFactory(),
		trEntity.Echo:   echoAdapter.NewFactory(),
	}

	// Config validators
	validators := map[trEntity.TaskRunnerType]runner.ConfigValidator{
		trEntity.AwsSqs: sqsAdapter.NewConfigValidator(),
		trEntity.Echo:   validation.NewConfigValidator([]runner.FieldRule{}),
	}

	// Task runner module
	taskRunnerRepository := trRepo.NewTaskRunnerRepository(dbConn)
	taskRunnerService := trSvc.NewTaskRunnerService(taskRunnerRepository, validators)
	dispatchService := dispatchSvc.NewDispatchService(taskRunnerRepository, factories)

	// Workflow module
	taskRepo := taskRepository.NewTaskRepository(dbConn)
	taskService := task.NewTaskService(taskRepo, dispatchService)

	workflowRepo := workflowRepository.NewWorkflowRepository(dbConn)
	wfService := workflowService.NewWorkflowService(workflowRepo)

	server := restapi.NewRestApiServer()
	server.RegisterRoutes(taskService, wfService, taskRunnerService)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
