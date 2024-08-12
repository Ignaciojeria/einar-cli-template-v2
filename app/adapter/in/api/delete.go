package api

import (
	"archetype/app/shared/infrastructure/serverwrapper"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/labstack/echo/v4"
)

func init() {
	ioc.Registry(newTemplateDelete, serverwrapper.NewEchoWrapper)
}
func newTemplateDelete(e serverwrapper.EchoWrapper) {
	e.DELETE("/insert-your-custom-pattern-here", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Unimplemented",
		})
	})
}
