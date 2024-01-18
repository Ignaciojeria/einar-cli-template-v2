package business

import (
	"context"
	"einar/app/domain"
	"einar/app/shared/container"

	"einar/app/domain/port/in"
)

var Example = container.InjectBusiness[in.Example](func() (in.Example, error) {
	instance := exampleUsecase{
		Hello: Other.Dependency,
	}
	return instance.run, nil
})

type exampleUsecase struct {
	DependencyID string
	Hello        in.Other
}

func (u exampleUsecase) run(ctx context.Context, e domain.Example) (string, error) {
	u.Hello(ctx, e)
	return "hello :D", nil
}
