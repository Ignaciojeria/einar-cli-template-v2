package business

import (
	"context"
	"einar/app/container"
	"einar/app/domain"

	"einar/app/domain/port/in"
)

var Example = container.InjectUseCase[in.Example](func() in.Example {
	instance := exampleUsecase{
		Hello: Other,
	}
	return instance.run
})

type exampleUsecase struct {
	DependencyID string
	Hello        in.Other
}

func (u exampleUsecase) run(ctx context.Context, e domain.Example) string {
	u.Hello(ctx, e)
	return "hello :D"
}
