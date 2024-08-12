package usecase

import (
	"context"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type INewUsecase func(ctx context.Context, input interface{}) (interface{}, error)

func init() {
	ioc.Registry(NewUseCase)
}

func NewUseCase() INewUsecase {
	return func(ctx context.Context, input interface{}) (interface{}, error) {
		return input, nil
	}
}
