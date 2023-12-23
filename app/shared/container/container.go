package container

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Container[T any] struct {
	loadableContent func() (T, error)
	Dependency      T
	Err             error
}

func (c *Container[T]) Load() (any, error) {
	instance, err := c.loadableContent()
	c.Dependency = instance
	c.Err = err
	return instance, err
}

type Loadable[T any] interface {
	Load() (any, error)
}

var Installations = make(map[string]Loadable[any])

func InjectInstallation[T any](loadableContent func() (T, error)) *Container[T] {
	adapter := Container[T]{loadableContent: loadableContent}
	Installations[uuid.NewString()] = &adapter
	return &adapter
}

var Business = make(map[string]Loadable[any])

func InjectBusiness[T any](loadableContent func() (T, error)) *Container[T] {
	adapter := Container[T]{loadableContent: loadableContent}
	Business[uuid.NewString()] = &adapter
	return &adapter
}

var InboundAdapters = make(map[string]Loadable[any])

func InjectInboundAdapter[T any](loadableContent func() (T, error)) Loadable[T] {
	adapter := Container[T]{loadableContent: loadableContent}
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
