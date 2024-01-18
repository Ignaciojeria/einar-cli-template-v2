package server

import (
	"einar/app/shared/config"
	"einar/app/shared/container"
	"fmt"
	"log"
	"sync"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type startHTTPServer func() error

var c = container.InjectInstallation[*echo.Echo](func() (*echo.Echo, error) {
	e := echo.New()
	e.Use(otelecho.Middleware(config.PROJECT_NAME.Get() + "-http-server"))
	return e, nil
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
	if c.Dependency == nil {
		return echo.New()
	}
	return c.Dependency
}
