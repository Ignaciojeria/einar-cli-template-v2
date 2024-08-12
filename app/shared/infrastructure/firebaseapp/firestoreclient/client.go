package firestoreclient

import (
	"archetype/app/shared/infrastructure/firebaseapp"

	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	ioc.Registry(NewClient, firebaseapp.NewFirebaseAPP)
}
func NewClient(app *firebase.App) (*firestore.Client, error) {
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	_, err = client.Collection("health").Doc("ping").Get(ctx)
	if status.Code(err) != codes.NotFound {
		return nil, err
	}
	return client, nil
}
