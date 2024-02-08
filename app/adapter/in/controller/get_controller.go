package controller

import (
	"archetype/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newGetController,
	server.NewServer)

type getController struct{}

func newGetController(s server.Server) getController {
	controller := getController{}
	s.Router().GET(s.ApiPrefix()+"insert_your_pattern", controller.handle)
	return controller
}

func (ctrl getController) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
