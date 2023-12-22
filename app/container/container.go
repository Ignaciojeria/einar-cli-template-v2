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

type Dependency[T any] func() T

var (
	ConfigurationContainer   = make(map[string]DependencyContainer[any])
	UseCaseContainer         = make(map[string]DependencyContainer[any])
	HTTPServerContainer      = make(map[string]DependencyContainer[any])
	InstallationContainer    = make(map[string]DependencyContainer[any])
	OutboundAdapterContainer = make(map[string]DependencyContainer[Adapter])
)

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

func InjectOutBoundAdapter[T any](t Dependency[Adapter]) Adapter {
	OutboundAdapterContainer[uuid.NewString()] = DependencyContainer[Adapter]{Dependency: t()}
	return t()
}
