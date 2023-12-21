package container

import (
	"github.com/google/uuid"
)

type Adapter func() error

func (a Adapter) Load() error {
	return a()
}

type DependencyContainer[T any] struct {
	Dependency T
}

var (
	ConfigurationContainer   = make(map[string]DependencyContainer[any])
	UseCaseContainer         = make(map[string]DependencyContainer[any])
	HTTPServerContainer      = make(map[string]DependencyContainer[any])
	InstallationContainer    = make(map[string]DependencyContainer[any])
	InboundAdapterContainer  = make(map[string]DependencyContainer[Adapter])
	OutboundAdapterContainer = make(map[string]DependencyContainer[Adapter])
)

type Dependency[T any] func() T

func InjectUseCase[T any](t Dependency[T]) T {
	UseCaseContainer[uuid.NewString()] = DependencyContainer[any]{Dependency: t()}
	return t()
}

func InjectInstallation[T any](t Dependency[T]) T {
	InstallationContainer[uuid.NewString()] = DependencyContainer[any]{Dependency: t()}
	return t()
}

func InjectHTTPServer[T any](t Dependency[T]) T {
	HTTPServerContainer[uuid.NewString()] = DependencyContainer[any]{Dependency: t()}
	return t()
}

func InjectInboundAdapter[T any](t Dependency[Adapter]) Adapter {
	InboundAdapterContainer[uuid.NewString()] = DependencyContainer[Adapter]{Dependency: t()}
	return t()
}

func InjectOutBoundAdapter[T any](t Dependency[Adapter]) Adapter {
	OutboundAdapterContainer[uuid.NewString()] = DependencyContainer[Adapter]{Dependency: t()}
	return t()
}
