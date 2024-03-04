package topic

import (
	"archetype/app/infrastructure/eventbroker"
	"sync"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc"
)

var _ = ioc.Registry(NewTopics, eventbroker.NewClient)

type Topics interface {
	Get(topicName string) *pubsub.Topic
}

type topics struct {
	topicRefs sync.Map
	client    *pubsub.Client
}

func NewTopics(c *pubsub.Client) topics {
	return topics{client: c}
}

func (t *topics) Get(topicName string) *pubsub.Topic {
	value, ok := t.topicRefs.Load(topicName)
	if ok {
		return value.(*pubsub.Topic)
	}
	newTopicRef := t.client.Topic(topicName)
	t.topicRefs.Store(topicName, newTopicRef)
	return newTopicRef
}
