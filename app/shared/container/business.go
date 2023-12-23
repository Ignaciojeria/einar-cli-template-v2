package container

import "github.com/google/uuid"

var Business = make(map[string]Loadable[any])

func InjectBusiness[T any](loadableContent func() (T, error)) *Container[T] {
	adapter := Container[T]{loadableContent: loadableContent}
	Business[uuid.NewString()] = &adapter
	return &adapter
}
