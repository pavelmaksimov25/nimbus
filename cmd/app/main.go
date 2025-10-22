package main

import (
	restserver "shiva/internal/presenter/http"
)

func main() {
	server := restserver.NewTaskWebServer()
	server.RegisterRoutes()
	err := server.Run()
	if err != nil {
		panic(err)
	}
}