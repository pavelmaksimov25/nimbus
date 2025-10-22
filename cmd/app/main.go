package main

import (
	restserver "shiva/internal/presenter/http"
)

func main() {
	server := restserver.NewTaskWebServer()
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}