package main

import (
	_ "einar/app/adapter/in/controller"
	_ "einar/app/adapter/in/subscription"
	"einar/app/infrastructure/server"
	"einar/app/shared/container"
	"os"
)

func main() {
	if err := container.LoadDependencies(); err != nil {
		os.Exit(0)
	}
	server.StartHTTPServer()
}
