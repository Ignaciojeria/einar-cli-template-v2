package natsconsumer

import (
	"archetype/app/shared/infrastructure/natsconn"
	"archetype/app/shared/logging"
	"context"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
)

func init() {
	ioc.Registry(newNatsConsumer, natsconn.NewJetStream, logging.NewLogger)
}
func newNatsConsumer(js jetstream.JetStream) (jetstream.ConsumeContext, error) {
	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, "stream-name", jetstream.ConsumerConfig{
		Name:          "consumer-name",
		Durable:       "consumer-name",
		MaxAckPending: 5,
	})
	if err != nil {
		return nil, err
	}
	return consumer.Consume(func(msg jetstream.Msg) {
		//process received message here
		msg.Ack()
	})
}
