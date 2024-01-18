package subscription

import (
	"context"
	"einar/app/shared/constants"
	internal "einar/app/shared/errors"
	"einar/app/shared/slog"
	"errors"
	"net/http"

	"cloud.google.com/go/pubsub"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type HandleMessageAcknowledgementDetails struct {
	SubscriptionName    string
	Error               error
	Message             *pubsub.Message
	MessageID           string
	PublishTime         string
	StatusCode          int
	ErrorsRequiringNack []error
	CustomLogFields     map[string]interface{}
}

func HandleMessageAcknowledgement(ctx context.Context, details *HandleMessageAcknowledgementDetails) int {
	ctx, span := tracer.Start(ctx,
		"HandleMessageAcknowledgement",
		trace.WithSpanKind(trace.SpanKindConsumer), trace.WithAttributes(
			attribute.String("subscription.name", details.SubscriptionName),
			attribute.String("message.id", details.MessageID),
			attribute.String("message.publishTime", details.PublishTime),
		))
	defer span.End()

	if details.Error != nil {
		span.RecordError(details.Error)
		span.SetStatus(codes.Error, details.Error.Error())
		slog.SpanLogger(span).Error(
			details.SubscriptionName+"_exception",
			subscription_name, details.SubscriptionName,
			constants.Fields, details.CustomLogFields,
			constants.Error, details.Error,
		)

		for _, err := range details.ErrorsRequiringNack {
			if errors.Is(details.Error, err) {
				span.AddEvent("Event processing nacked, retrying",
					trace.WithAttributes(attribute.String(constants.Error, details.Error.Error())))
				details.Message.Nack()
				if errors.Is(err, internal.INTERNAL_SERVER_ERROR) {
					return http.StatusInternalServerError
				}
				return http.StatusBadGateway
			}
		}

		span.AddEvent("Event discarded",
			trace.WithAttributes(attribute.String(constants.Error, details.Error.Error())))
		details.Message.Ack()

		if details.Message.Attributes[constants.EventChannel] == constants.GoogleCloud {
			return http.StatusNoContent
		}

		return details.StatusCode
	}

	span.SetStatus(codes.Ok, details.SubscriptionName+"_succedded")

	slog.SpanLogger(span).Info(
		details.SubscriptionName+"_succedded",
		subscription_name, details.SubscriptionName,
		constants.Fields, details.CustomLogFields,
	)
	details.Message.Ack()

	return details.StatusCode
}
