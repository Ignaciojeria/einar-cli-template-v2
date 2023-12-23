package controller

import (
	"context"
	"einar/app/domain"
	"einar/app/domain/port/in"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetExampleControllerHandle(t *testing.T) {
	// Obtener la instancia de getExampleController
	instance, err := getExampleInstance.Load()
	assert.NoError(t, err)

	// Type assertion para obtener la instancia de getExampleController
	controller, ok := instance.(getExampleController)
	assert.True(t, ok)

	// Obtener el patrón de la instancia
	pattern := controller.pattern

	// Mock de Echo y Example
	e := echo.New()
	var mockExample in.Example = func(ctx context.Context, e domain.Example) (string, error) {
		return "Hello mom", nil
	}

	// Crear instancia del controlador con mock
	controller = getExampleController{
		example: mockExample,
		pattern: pattern,
	}

	// Registrar el endpoint en Echo con el patrón obtenido
	e.GET(pattern, controller.handle)

	// Crear una solicitud HTTP de prueba utilizando el mismo patrón
	req := httptest.NewRequest(http.MethodGet, pattern, nil)
	rec := httptest.NewRecorder()

	// Crear un contexto Echo y ejecutar el handler
	c := e.NewContext(req, rec)
	err = controller.handle(c)

	// Verificar resultados
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	// Aquí puedes añadir más assertions para verificar el body de la respuesta, headers, etc.
}
