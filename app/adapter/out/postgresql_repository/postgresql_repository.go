package postgresql_repository

import (
	"archetype/app/shared/infrastructure/postgresql"
	"context"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type RunPostgreSQLOperation func(ctx context.Context, input interface{}) error

func init() {
	ioc.Registry(
		NewRunPostgreSQLOperation,
		postgresql.NewConnection)
}
func NewRunPostgreSQLOperation(connection *gorm.DB) RunPostgreSQLOperation {
	return func(ctx context.Context, input interface{}) error {
		return nil
	}
}
