package configuration

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type StorjConfiguration struct {
	STORJ_ACCESS_GRANT string `required:"true"`
}

func init() {
	ioc.Registry(NewStorjConfiguration, NewEnvLoader)
}
func NewStorjConfiguration(env EnvLoader) (StorjConfiguration, error) {
	conf := StorjConfiguration{
		STORJ_ACCESS_GRANT: env.Get("STORJ_ACCESS_GRANT"),
	}
	return validateConfig(conf)
}
