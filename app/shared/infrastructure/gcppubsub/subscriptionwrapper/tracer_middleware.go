package subscriptionwrapper

import (
	"archetype/app/shared/logging"
	"context"

	"cloud.google.com/go/pubsub"
)

// TODO : CREAR INSTANCIA INTERNA DE LOGGER Y QUITARLO DE LA INYECCIÓN DE LA DEPENDENCIA.
func NewTracerMiddleware(logger logging.Logger) Middleware {
	return func(next MessageProcessor) MessageProcessor {
		return func(ctx context.Context, m *pubsub.Message) (int, error) {
			logger.Info("Processing message with ID: " + m.ID)
			status, err := next(ctx, m)
			if err != nil {
				logger.Error("Error processing message with ID: "+m.ID, "error", err.Error())
			} else {
				logger.Info("Successfully processed message with ID: " + m.ID)
			}
			return status, err
		}
	}
}
