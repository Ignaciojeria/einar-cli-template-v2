package controller

import (
	"archetype/app/configuration"
	"archetype/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newPutController,
	server.NewServer,
	configuration.NewConf)

type putController struct {
	s server.Server
}

func newPutController(
	s server.Server,
	c configuration.Conf) putController {
	controller := putController{
		s: s,
	}
	controller.s.Router().PUT(c.ApiPrefix+"insert_your_pattern", controller.handle)
	return controller
}

func (ctrl putController) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
