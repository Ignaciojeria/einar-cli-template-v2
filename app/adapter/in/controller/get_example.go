package controller

import (
	"einar/app/business"
	"einar/app/container"
	"einar/app/domain"
	"einar/app/domain/port/in"
	"einar/app/infrastructure/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

var getExampleInstance = container.InjectInboundAdapter(func() (getExampleController, error) {
	instance := getExampleController{
		echo:    server.Echo,
		example: business.Example,
		pattern: "/api/pattern",
	}
	instance.echo.GET(instance.pattern, instance.handle)
	return instance, nil
})

type getExampleController struct {
	echo    *echo.Echo
	pattern string
	example in.Example
}

func (u getExampleController) handle(c echo.Context) error {
	exampleResponse := u.example(c.Request().Context(), domain.Example{})
	return c.JSON(http.StatusOK, exampleResponse)
}
