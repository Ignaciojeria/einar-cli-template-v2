package container

import "github.com/google/uuid"

var InboundAdapters = make(map[string]Loadable[any])

func InjectInboundAdapter[T any](loadableContent func() (T, error)) Loadable[T] {
	adapter := Container[T]{loadableContent: loadableContent}
	InboundAdapters[uuid.NewString()] = &adapter
	return &adapter
}
