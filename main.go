package main

import (
	_ "einar/app/adapter/in/controller"
	"einar/app/container"
	"einar/app/infrastructure/server"
)

func main() {
	for _, v := range container.InboundAdapterContainer {
		v.Load()
	}
	server.StartHTTPServer()
}
