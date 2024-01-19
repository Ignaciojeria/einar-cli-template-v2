package controller

import (
	"archetype/app/configuration"
	"archetype/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newPostController,
	server.NewServer,
	configuration.NewConf)

type postController struct {
	s server.Server
}

func newPostController(
	s server.Server,
	c configuration.Conf) postController {
	controller := postController{
		s: s,
	}
	controller.s.Router().POST(c.ApiPrefix+"insert_your_pattern", controller.handle)
	return controller
}

func (ctrl postController) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
