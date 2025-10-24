package main

import (
	"nimbus/internal/workflow/adapters/store"
	restapi "nimbus/internal/workflow/presenter/rest_api"
)

func main() {
	store := store.NewTaskStoreInMemory()
	server := restapi.NewRestApiServer()
	server.RegisterRoutes(store)// it breaks layers. Only for demo purposes
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
