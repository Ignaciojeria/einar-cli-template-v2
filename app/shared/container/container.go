package container

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
