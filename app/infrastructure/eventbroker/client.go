package eventbroker

import (
	"archetype/app/configuration"
	"context"

	"log/slog"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc"
)

var _ = ioc.Registry(NewClient, configuration.NewConf)

const DefaultPushHandlerPrefix = "/topic/"

func NewClient(conf configuration.Conf) (*pubsub.Client, error) {
	c, err := pubsub.NewClient(context.Background(), conf.GoogleProjectID)
	if err != nil {
		slog.Error("error getting pubsub client", err.Error())
		return nil, err
	}
	return c, nil
}
