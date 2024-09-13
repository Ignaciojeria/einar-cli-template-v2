package slogging

import (
	"archetype/app/shared/configuration"
	"log/slog"
	"os"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel/trace"
)

const (
	envKey     = "env"
	serviceKey = "service"
	traceIDKey = "trace_id"
	spanIDKey  = "span_id"
)

type SpanLogger struct {
	config configuration.Conf
	*slog.Logger
}

func (logger *SpanLogger) SpanLogger(span trace.Span) *slog.Logger {
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()
	return logger.With(
		slog.String(envKey, logger.config.ENVIRONMENT),
		slog.String(serviceKey, logger.config.PROJECT_NAME),
		slog.String(traceIDKey, traceID),
		slog.String(spanIDKey, spanID),
	)
}

func init() {
	ioc.Registry(NewSpanLogger, configuration.NewConf)
}

func NewSpanLogger(conf configuration.Conf) SpanLogger {
	return SpanLogger{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		config: conf,
	}
}
