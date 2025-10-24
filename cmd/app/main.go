package main

import (
	"nimbus/internal/workflow/adapters/store"
	"nimbus/internal/workflow/application/task"
	restapi "nimbus/internal/workflow/presenter/rest_api"
)

func main() {
	taskStore := store.NewTaskStoreInMemory()
	taskService := task.NewTaskService(taskStore)
	server := restapi.NewRestApiServer()
	server.RegisterRoutes(taskService)
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
