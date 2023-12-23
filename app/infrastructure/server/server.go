package server

import (
	"einar/app/shared/container"
	"fmt"
	"log"
	"sync"

	"github.com/labstack/echo/v4"
)

type startHTTPServer func() error

var c = container.InjectInstallation[*echo.Echo](func() (*echo.Echo, error) {
	return echo.New(), nil
})

var once sync.Once
var StartHTTPServer = func() {
	for _, route := range Echo().Routes() {
		fmt.Printf("Method: %v, Path: %v, Name: %v\n", route.Method, route.Path, route.Name)
	}
	err := Echo().Start(":" + "8080")
	if err != nil {
		log.Panic(err)
	}
}

func Echo() *echo.Echo {
	return c.Dependency
}
