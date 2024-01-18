package server

import (
	"archetype/app/configuration"
	"fmt"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(NewServer, configuration.NewConf)

type Server struct {
	c configuration.Conf
	e *echo.Echo
}

func NewServer(c configuration.Conf) Server {
	e := echo.New()
	return Server{
		e: e,
		c: c,
	}
}

func (s Server) Start() {
	s.printRoutes()
	s.e.Start(":" + s.c.Port)
}

func (s Server) printRoutes() {
	for _, route := range s.e.Routes() {
		fmt.Printf("Method: %v, Path: %v, Name: %v\n", route.Method, route.Path, route.Name)
	}
}

func (s Server) Router() *echo.Echo {
	return s.e
}
