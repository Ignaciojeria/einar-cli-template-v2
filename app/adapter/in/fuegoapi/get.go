package fuegoapi

import (
	"archetype/app/shared/infrastructure/fuegoserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(newTemplateGet, fuegoserver.New)
}
func newTemplateGet(s *fuego.Server) {
	fuego.Get(s, "/insert-your-custom-pattern-here", func(c *fuego.ContextWithBody[any]) (any, error) {
		body, err := c.Body()
		if err != nil {
			return "unimplemented", err
		}
		return body, nil
	}, option.Summary("newTemplateDelete"))
}
