package main

import (
	_ "einar/app/adapter/in/controller"
	_ "einar/app/adapter/in/subscription"
	"einar/app/infrastructure/server"
	"einar/app/shared/container"
	"os"

	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found getting environments from system")
	}
	if err := container.LoadDependencies(); err != nil {
		os.Exit(0)
	}
	server.StartHTTPServer()
}
