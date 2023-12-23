package business

import (
	"context"
	"einar/app/domain"
	"einar/app/domain/port/in"
	"einar/app/shared/container"
)

var Other = container.InjectUseCase[in.Other](func() in.Other {
	instance := exampleOther{}
	return instance.run
})

type exampleOther struct {
}

func (u exampleOther) run(ctx context.Context, e domain.Example) string {
	return "hello :D"
}
