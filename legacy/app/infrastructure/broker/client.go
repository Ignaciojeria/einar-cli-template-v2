package broker

import (
	"context"
	"einar/app/shared/config"
	"einar/app/shared/constants"
	"einar/app/shared/container"
	"einar/app/shared/slog"

	"cloud.google.com/go/pubsub"
)

var c = container.InjectInstallation[*pubsub.Client](func() (*pubsub.Client, error) {
	config.Installations.EnablePubSub = true
	projectId := config.GOOGLE_PROJECT_ID.Get()
	c, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		slog.Logger().Error("error getting pubsub client", constants.Error, err.Error())
		return c, err
	}
	return c, nil
})

func Client() *pubsub.Client {
	if c.Dependency == nil {
		return &pubsub.Client{}
	}
	return c.Dependency
}
