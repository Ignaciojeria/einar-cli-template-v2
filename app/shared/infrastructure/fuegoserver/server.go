package fuegoserver

import (
	"archetype/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/hellofresh/health-go/v5"
)

func init() {
	ioc.Registry(New)
	ioc.Registry(healthCheck, New, configuration.NewConf)
	ioc.RegistryAtEnd(startAtEnd, New)
}

func New() *fuego.Server {
	return fuego.NewServer()
}

func startAtEnd(e *fuego.Server) error {
	return e.Run()
}

func healthCheck(s *fuego.Server, c configuration.Conf) error {
	h, err := health.New(
		health.WithComponent(health.Component{
			Name:    c.PROJECT_NAME,
			Version: c.VERSION,
		}), health.WithSystemInfo())
	if err != nil {
		return err
	}
	fuego.GetStd(s,
		"/health",
		h.Handler().ServeHTTP,
		option.Summary("healthCheck"))
	return nil
}
