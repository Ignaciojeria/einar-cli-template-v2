package slogging

import (
	"archetype/app/shared/configuration"
	"log/slog"
	"os"
	"strconv"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel/trace"
)

// Datadog trace and log correlation :
// https://docs.datadoghq.com/tracing/other_telemetry/connect_logs_and_traces/opentelemetry/?tab=go
const (
	ddTraceIDKey = "dd.trace_id"
	ddSpanIDKey  = "dd.span_id"
	ddServiceKey = "dd.service"
	ddEnvKey     = "dd.env"
	ddVersionKey = "dd.version"
)

type DDSpanLogger struct {
	config configuration.Conf
	*slog.Logger
}

func (logger *DDSpanLogger) SpanLogger(span trace.Span) *slog.Logger {
	const (
		envKey     = "env"
		serviceKey = "service"
		traceIDKey = "trace_id"
		spanIDKey  = "span_id"
	)
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	ddService := logger.config.DD_SERVICE
	ddEnv := logger.config.DD_ENV
	ddVersion := logger.config.DD_VERSION

	if ddService == "" || ddEnv == "" || ddVersion == "" {
		return logger.With(
			slog.String(envKey, logger.config.ENVIRONMENT),
			slog.String(serviceKey, logger.config.PROJECT_NAME),
			slog.String(traceIDKey, traceID),
			slog.String(spanIDKey, spanID),
		)
	}
	return logger.With(
		slog.String(traceIDKey, traceID),
		slog.String(spanIDKey, spanID),
		slog.String(ddTraceIDKey, convertTraceID(traceID)),
		slog.String(ddSpanIDKey, convertTraceID(spanID)),
		slog.String(ddServiceKey, ddService),
		slog.String(ddEnvKey, ddEnv),
		slog.String(ddVersionKey, ddVersion),
	)
}

func init() {
	ioc.Registry(NewDDSpanLogger, configuration.NewConf)
}
func NewDDSpanLogger(conf configuration.Conf) DDSpanLogger {
	return DDSpanLogger{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		config: conf,
	}
}

func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
