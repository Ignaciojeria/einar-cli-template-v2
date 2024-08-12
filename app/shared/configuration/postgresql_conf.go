package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type PostgreSQLConfiguration struct {
	DATABASE_POSTGRES_HOSTNAME string `required:"true"`
	DATABASE_POSTGRES_PORT     string `required:"true"`
	DATABASE_POSTGRES_NAME     string `required:"true"`
	DATABASE_POSTGRES_USERNAME string `required:"true"`
	DATABASE_POSTGRES_PASSWORD string `required:"true"`
	DATABASE_POSTGRES_SSL_MODE string `required:"true"`
}

func init() {
	ioc.Registry(NewPostgreSQLConfiguration, NewEnvLoader)
}
func NewPostgreSQLConfiguration(env EnvLoader) (PostgreSQLConfiguration, error) {
	conf := PostgreSQLConfiguration{
		DATABASE_POSTGRES_HOSTNAME: env.Get("DATABASE_POSTGRES_HOSTNAME"),
		DATABASE_POSTGRES_PORT:     env.Get("DATABASE_POSTGRES_PORT"),
		DATABASE_POSTGRES_NAME:     env.Get("DATABASE_POSTGRES_NAME"),
		DATABASE_POSTGRES_USERNAME: env.Get("DATABASE_POSTGRES_USERNAME"),
		DATABASE_POSTGRES_PASSWORD: env.Get("DATABASE_POSTGRES_PASSWORD"),
		DATABASE_POSTGRES_SSL_MODE: env.Get("DATABASE_POSTGRES_SSL_MODE"),
	}
	return validateConfig(conf)
}
