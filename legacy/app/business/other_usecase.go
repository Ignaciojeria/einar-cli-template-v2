package business

import (
	"context"
	"einar/app/domain"
	"einar/app/domain/port/in"
	"einar/app/shared/container"
)

var Other = container.InjectBusiness[in.Other](func() (in.Other, error) {
	instance := exampleOther{}
	return instance.run, nil
})

type exampleOther struct {
}

func (u exampleOther) run(ctx context.Context, e domain.Example) string {
	return "hello :D"
}
