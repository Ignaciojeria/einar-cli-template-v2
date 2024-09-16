package gcppublisher

import (
	"archetype/app/shared/infrastructure/gcppubsub"
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type INewPublishEvent func(ctx context.Context, input interface{}) error

func init() {
	ioc.Registry(
		NewPublishEvent,
		gcppubsub.NewClient)
}
func NewPublishEvent(c *pubsub.Client) INewPublishEvent {
	topicName := "INSERT_YOUR_TOPIC_NAME_HERE"
	topic := c.Topic(topicName)
	return func(ctx context.Context, input interface{}) error {

		bytes, err := json.Marshal(input)
		if err != nil {
			return err
		}

		message := &pubsub.Message{
			Attributes: map[string]string{
				"customAttribute1": "attr1",
				"customAttribute2": "attr2",
			},
			Data: bytes,
		}

		result := topic.Publish(ctx, message)
		// Get the server-generated message ID.
		_, err = result.Get(ctx)

		if err != nil {
			return err
		}

		return nil
	}
}
