package server

import (
	"einar/app/container"
	"fmt"
	"log"
	"sync"

	"github.com/labstack/echo/v4"
)

type startHTTPServer func() error

var Echo = container.InjectInstallation[*echo.Echo](func() *echo.Echo {
	return echo.New()
})

var once sync.Once
var StartHTTPServer = container.InjectHTTPServer[startHTTPServer](func() startHTTPServer {
	return func() error {
		once.Do(func() {
			for _, route := range Echo.Routes() {
				fmt.Printf("Method: %v, Path: %v, Name: %v\n", route.Method, route.Path, route.Name)
			}
			err := Echo.Start(":" + "8080")
			if err != nil {
				log.Panic(err)
			}
		})
		return nil
	}
})
