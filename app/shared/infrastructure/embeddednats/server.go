package embeddednats

import (
	"archetype/app/shared/configuration"
	"errors"
	"net/url"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func init() {
	ioc.Registry(NewConn, configuration.NewNatsConfiguration)
}
func NewConn(conf configuration.NatsConfiguration) (*nats.Conn, error) {
	leafURL, err := url.Parse("nats-leaf://" + conf.NATS_CONNECTION_URL)
	if err != nil {
		return nil, err
	}
	opts := server.Options{
		ServerName:      "embedded_server",
		DontListen:      true,
		JetStream:       true,
		JetStreamDomain: "embedded",
		LeafNode: server.LeafNodeOpts{
			Remotes: []*server.RemoteLeafOpts{
				{
					URLs:        []*url.URL{leafURL},
					Credentials: conf.NATS_CONNECTION_CREDS_FILEPATH,
				},
			},
		},
	}
	ns, err := server.NewServer(&opts)
	if err != nil {
		return nil, err
	}
	ns.ConfigureLogger()
	go ns.Start()
	if !ns.ReadyForConnections(5 * time.Second) {
		return nil, errors.New("NATS Server timeout")
	}
	clientOpts := []nats.Option{
		nats.InProcessServer(ns),
	}
	conn, err := nats.Connect(ns.ClientURL(), clientOpts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
