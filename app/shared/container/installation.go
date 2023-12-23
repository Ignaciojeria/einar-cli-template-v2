package container

import "github.com/google/uuid"

var Installations = make(map[string]Loadable[any])

func InjectInstallation[T any](loadableContent func() (T, error)) *Container[T] {
	adapter := Container[T]{loadableContent: loadableContent}
	Installations[uuid.NewString()] = &adapter
	return &adapter
}
