package firebaseapp

import (
	"archetype/app/shared/configuration"

	"context"

	firebase "firebase.google.com/go"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewFirebaseAPP, configuration.NewConf)
}
func NewFirebaseAPP(conf configuration.Conf) (*firebase.App, error) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: conf.GOOGLE_PROJECT_ID,
	})
	if err != nil {
		return nil, err
	}
	return app, nil
}
