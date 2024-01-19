package controller

import (
	"archetype/app/configuration"
	"archetype/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newPatchController,
	server.NewServer,
	configuration.NewConf)

type patchController struct {
	s server.Server
}

func newPatchController(
	s server.Server,
	c configuration.Conf) patchController {
	controller := patchController{
		s: s,
	}
	controller.s.Router().PATCH(c.ApiPrefix+"insert_your_pattern", controller.handle)
	return controller
}

func (ctrl patchController) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
