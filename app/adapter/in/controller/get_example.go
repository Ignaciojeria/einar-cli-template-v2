package controller

import (
	"einar/app/business"
	"einar/app/domain"
	"einar/app/domain/port/in"
	"einar/app/infrastructure/server"
	"einar/app/shared/container"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getExampleController struct {
	pattern string
	example in.Example
}

var getExampleInstance = container.InjectInboundAdapter(func() (getExampleController, error) {
	instance := getExampleController{
		example: business.Example.Dependency,
		pattern: "/api/pattern",
	}
	server.Echo().GET(instance.pattern, instance.handle)
	return instance, nil
})

func (u getExampleController) handle(c echo.Context) error {
	exampleResponse, _ := u.example(c.Request().Context(), domain.Example{})
	return c.JSON(http.StatusOK, exampleResponse)
}
