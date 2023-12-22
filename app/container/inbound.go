package container

import "github.com/google/uuid"

// AdapterFactory envuelve una función de fábrica para crear instancias de tipo T.
type AdapterFactory[T any] struct {
	factory func() (T, error)
}

func (af AdapterFactory[T]) Load() (any, error) {
	instance, err := af.factory()
	return instance, err // Aquí, 'instance' es automáticamente considerado 'any'
}

// Loadable define una interfaz para objetos que pueden ser cargados.
type Loadable interface {
	Load() (any, error)
}

// InboundAdapterContainer almacena adaptadores que pueden ser cargados.
var InboundAdapterContainer = make(map[string]Loadable)

// InjectInboundAdapter inyecta un adaptador en el contenedor y devuelve un Loadable.
// El identificador proporcionado se usa como clave en el contenedor.
func InjectInboundAdapter[T any](factory func() (T, error)) Loadable {
	adapter := AdapterFactory[T]{factory: factory}
	InboundAdapterContainer[uuid.NewString()] = adapter
	return adapter
}
