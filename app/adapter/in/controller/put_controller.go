package controller

import (
	"archetype/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newPutController,
	server.NewServer)

type putController struct{}

func newPutController(s server.Server) putController {
	controller := putController{}
	s.Router().PUT(s.ApiPrefix()+"insert_your_pattern", controller.handle)
	return controller
}

func (ctrl putController) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
