package postgresql

import (
	"archetype/app/shared/configuration"
	"fmt"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func init() {
	ioc.Registry(
		NewConnection,
		configuration.NewPostgreSQLConfiguration)
}
func NewConnection(env configuration.PostgreSQLConfiguration) (*gorm.DB, error) {
	username := env.DATABASE_POSTGRES_USERNAME
	pwd := env.DATABASE_POSTGRES_PASSWORD
	host := env.DATABASE_POSTGRES_HOSTNAME
	port := env.DATABASE_POSTGRES_PORT
	dbname := env.DATABASE_POSTGRES_NAME
	sslMode := env.DATABASE_POSTGRES_SSL_MODE
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", username, pwd, host, port, dbname, sslMode)))
	if err != nil {
		return nil, err
	}
	if err := db.Use(tracing.NewPlugin()); err != nil {
		return nil, err
	}
	return db, nil
}
