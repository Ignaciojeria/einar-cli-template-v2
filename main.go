package main

import (
	_ "einar/app/adapter/in/controller"
	_ "einar/app/adapter/in/subscription"
	"einar/app/infrastructure/server"
	"einar/app/shared/container"

	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found getting environments from system")
	}
	for _, v := range container.Installations {
		v.Load()
	}
	for _, v := range container.InboundAdapters {
		v.Load()
	}
	server.StartHTTPServer()
}
