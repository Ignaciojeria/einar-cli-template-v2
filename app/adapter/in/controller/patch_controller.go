package controller

import (
	"archetype/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newPatchController,
	server.NewServer)

type patchController struct{}

func newPatchController(s server.Server) patchController {
	controller := patchController{}
	s.Router().PATCH(s.ApiPrefix()+"insert_your_pattern", controller.handle)
	return controller
}

func (ctrl patchController) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
