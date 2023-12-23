package in

import (
	"context"
	"einar/app/domain"
)

type Example func(ctx context.Context, e domain.Example) (string, error)
