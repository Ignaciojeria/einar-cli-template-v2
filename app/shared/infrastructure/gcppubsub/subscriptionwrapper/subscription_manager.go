package subscriptionwrapper

import (
	"archetype/app/shared/infrastructure/gcppubsub"
	"archetype/app/shared/infrastructure/serverwrapper"
	"log/slog"
	"os"

	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/labstack/echo/v4"
)

type SubscriptionManager interface {
	Subscription(id string) *pubsub.Subscription
	WithMessageProcessor(mp MessageProcessor) SubscriptionManager
	WithPushHandler(path string) SubscriptionManager
	Start(subscriptionRef *pubsub.Subscription) (SubscriptionManager, error)
}

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type SubscriptionWrapper struct {
	client           *pubsub.Client
	httpServer       serverwrapper.EchoWrapper
	messageProcessor MessageProcessor
}

func init() {
	ioc.Registry(
		NewSubscriptionManager,
		gcppubsub.NewClient,
		serverwrapper.NewEchoWrapper,
	)
}
func NewSubscriptionManager(
	c *pubsub.Client,
	s serverwrapper.EchoWrapper) SubscriptionManager {
	return &SubscriptionWrapper{client: c, httpServer: s}
}

func newSubscriptionManagerWithMessageProcessor(
	c *pubsub.Client,
	s serverwrapper.EchoWrapper,
	mp MessageProcessor) SubscriptionManager {
	return &SubscriptionWrapper{client: c, httpServer: s, messageProcessor: mp}
}

func (sw *SubscriptionWrapper) Subscription(id string) *pubsub.Subscription {
	return sw.client.Subscription(id)
}

func (sw *SubscriptionWrapper) WithMessageProcessor(mp MessageProcessor) SubscriptionManager {
	return newSubscriptionManagerWithMessageProcessor(sw.client, sw.httpServer, mp)
}

func (s *SubscriptionWrapper) Start(subscriptionRef *pubsub.Subscription) (SubscriptionManager, error) {
	ctx := context.Background()
	if err := subscriptionRef.Receive(ctx, s.receive); err != nil {
		logger.Error(
			"subscription_signal_broken",
			"subscription_name", subscriptionRef.String(),
			"error", err.Error(),
		)
		time.Sleep(10 * time.Second)
		go s.Start(subscriptionRef)
		return s, err
	}
	return s, nil
}

func (s *SubscriptionWrapper) receive(ctx context.Context, m *pubsub.Message) {
	s.messageProcessor(ctx, m)
}

func (s *SubscriptionWrapper) WithPushHandler(path string) SubscriptionManager {
	s.httpServer.POST(path, s.pushHandler)
	return s
}

func (s *SubscriptionWrapper) pushHandler(c echo.Context) error {
	googleChannel := c.Request().Header.Get("X-Goog-Channel-ID")

	var msg pubsub.Message

	if googleChannel != "" {
		if err := c.Bind(&msg); err != nil {
			return c.String(http.StatusNoContent, "error binding Pub/Sub message")
		}
		statusCode, err := s.messageProcessor(c.Request().Context(), &msg)
		if statusCode >= 500 && statusCode <= 599 {
			return c.String(statusCode, "")
		}
		if err != nil {
			return c.String(http.StatusNoContent, "")
		}
		return c.String(http.StatusOK, "")
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "error reading request body")
	}

	msg.Attributes = make(map[string]string)
	for key, values := range c.Request().Header {
		if len(values) > 0 {
			msg.Attributes[strings.ToLower(key)] = strings.Join(values, ",")
		}
	}

	msg.Data = body
	if statusCode, err := s.messageProcessor(c.Request().Context(), &msg); err != nil {
		return c.String(statusCode, err.Error())
	}
	return c.String(http.StatusOK, "")
}
