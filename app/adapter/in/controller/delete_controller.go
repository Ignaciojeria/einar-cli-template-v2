package controller

import (
	"archetype/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newDeleteController,
	server.NewServer)

type deleteController struct{}

func newDeleteController(s server.Server) deleteController {
	controller := deleteController{}
	s.Router().DELETE(s.ApiPrefix()+"insert_your_pattern", controller.handle)
	return controller
}

func (ctrl deleteController) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
