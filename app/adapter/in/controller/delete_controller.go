package controller

import (
	"archetype/app/configuration"
	"archetype/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newDeleteController,
	server.NewServer,
	configuration.NewConf)

type deleteController struct {
	s server.Server
}

func newDeleteController(
	s server.Server,
	c configuration.Conf) deleteController {
	controller := deleteController{
		s: s,
	}
	controller.s.Router().DELETE(c.ApiPrefix+"insert_your_pattern", controller.handle)
	return controller
}

func (ctrl deleteController) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
