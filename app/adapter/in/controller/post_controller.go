package controller

import (
	"archetype/app/infrastructure/server"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/labstack/echo/v4"
)

var _ = ioc.Registry(
	newPostController,
	server.NewServer)

type postController struct{}

func newPostController(s server.Server) postController {
	controller := postController{}
	s.Router().POST(s.ApiPrefix()+"insert_your_pattern", controller.handle)
	return controller
}

func (ctrl postController) handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
