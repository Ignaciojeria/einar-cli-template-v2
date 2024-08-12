package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type NatsConfiguration struct {
	NATS_CONNECTION_URL            string `required:"false"`
	NATS_CONNECTION_CREDS_FILEPATH string `required:"true"`
}

func init() {
	ioc.Registry(NewNatsConfiguration, NewEnvLoader)
}
func NewNatsConfiguration(env EnvLoader) (NatsConfiguration, error) {
	conf := NatsConfiguration{
		NATS_CONNECTION_URL:            env.Get("NATS_CONNECTION_URL"),
		NATS_CONNECTION_CREDS_FILEPATH: env.Get("NATS_CONNECTION_CREDS_FILEPATH"),
	}
	return validateConfig(conf)
}
