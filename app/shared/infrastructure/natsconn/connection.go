package natsconn

import (
	"archetype/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go"
)

func init() {
	ioc.Registry(NewConn, configuration.NewNatsConfiguration)
}
func NewConn(conf configuration.NatsConfiguration) (*nats.Conn, error) {
	return nats.Connect(
		conf.NATS_CONNECTION_URL,
		nats.UserCredentials(conf.NATS_CONNECTION_CREDS_FILEPATH))
}
