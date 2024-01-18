package controller

import (
	"archetype/app/configuration"
	"archetype/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newGetController,
	server.NewServer,
	configuration.NewConf)

type getController struct {
	s server.Server
}

func newGetController(
	s server.Server,
	c configuration.Conf) getController {
	controller := getController{
		s: s,
	}
	controller.s.Router().GET(c.ApiPrefix+"insert_your_pattern", controller.handle)
	return controller
}

func (ctrl getController) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
