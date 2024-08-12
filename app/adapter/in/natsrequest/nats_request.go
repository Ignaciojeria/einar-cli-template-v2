package natsrequest

import (
	"archetype/app/shared/infrastructure/natsconn"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go"
)

func init() {
	ioc.Registry(newNatsRequest, natsconn.NewConn)
}
func newNatsRequest(nc *nats.Conn) (*nats.Subscription, error) {
	return nc.QueueSubscribe("example.request", "myQueueGroup", func(msg *nats.Msg) {
		nc.Publish(msg.Reply, []byte("example response"))
	})
}
