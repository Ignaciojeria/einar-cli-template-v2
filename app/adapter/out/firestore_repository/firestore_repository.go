package firestore_repository

import (
	"archetype/app/shared/infrastructure/firebaseapp/firestoreclient"
	"context"

	"cloud.google.com/go/firestore"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type IRunFirestoreOperation func(ctx context.Context, input interface{}) error

func init() {
	ioc.Registry(
		NewRunFirestoreOperation,
		firestoreclient.NewClient)
}
func NewRunFirestoreOperation(c *firestore.Client) IRunFirestoreOperation {
	return func(ctx context.Context, input interface{}) error {
		return nil
	}
}
