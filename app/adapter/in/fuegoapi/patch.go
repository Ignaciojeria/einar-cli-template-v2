package fuegoapi

import (
	"archetype/app/shared/infrastructure/fuegoserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
)

func init() {
	ioc.Registry(newTemplatePatch, fuegoserver.New)
}
func newTemplatePatch(s *fuego.Server) {
	fuego.Patch(s, "/insert-your-custom-pattern-here", func(c *fuego.ContextWithBody[any]) (any, error) {
		body, err := c.Body()
		if err != nil {
			return "unimplemented", err
		}
		return body, nil
	})
}
