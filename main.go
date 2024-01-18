package main

import (
	_ "archetype/app/adapter/in/controller"
	"archetype/app/infrastructure/server"
	"log"

	ioc "github.com/Ignaciojeria/einar-ioc"
)

func main() {
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
	s, _ := ioc.Get(server.NewServer)
	s.(server.Server).Start()
}
