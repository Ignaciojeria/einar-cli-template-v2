package postgresql

import (
	"archetype/app/shared/configuration"

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
	dbname := env.DATABASE_POSTGRES_NAME
	sslMode := env.DATABASE_POSTGRES_SSL_MODE
	db, err := gorm.Open(postgres.Open("postgres://" + username + ":" + pwd + "@" + host + "/" + dbname + "?sslmode=" + sslMode))
	if err != nil {
		return nil, err
	}
	if err := db.Use(tracing.NewPlugin()); err != nil {
		return nil, err
	}
	return db, nil
}
