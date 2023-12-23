package container

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Container[T any] struct {
	loadableDependency func() (T, error)
	isLoaded           bool
	Dependency         T
}

func (c *Container[T]) Load() (any, error) {
	if c.isLoaded {
		return nil, errors.New("dependency already loaded")
	}
	instance, err := c.loadableDependency()
	c.Dependency = instance
	c.isLoaded = true
	return instance, err
}

type Loadable[T any] interface {
	Load() (any, error)
}

var Installations = make(map[string]Loadable[any])

func InjectInstallation[T any](loadableDependency func() (T, error)) *Container[T] {
	adapter := Container[T]{loadableDependency: loadableDependency}
	Installations[uuid.NewString()] = &adapter
	return &adapter
}

var Business = make(map[string]Loadable[any])

func InjectBusiness[T any](loadableDependency func() (T, error)) *Container[T] {
	adapter := Container[T]{loadableDependency: loadableDependency}
	Business[uuid.NewString()] = &adapter
	return &adapter
}

var InboundAdapters = make(map[string]Loadable[any])

func InjectInboundAdapter[T any](loadableDependency func() (T, error)) Loadable[T] {
	adapter := Container[T]{loadableDependency: loadableDependency}
	InboundAdapters[uuid.NewString()] = &adapter
	return &adapter
}

func LoadDependencies() error {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found getting environments from system")
	}
	for _, v := range Installations {
		_, err := v.Load()
		if err != nil {
			return err
		}
	}
	for _, v := range Business {
		_, err := v.Load()
		if err != nil {
			return err
		}
	}
	for _, v := range InboundAdapters {
		_, err := v.Load()
		if err != nil {
			return err
		}
	}
	return nil
}
