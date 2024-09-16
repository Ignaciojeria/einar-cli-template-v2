package observability

import (
	"log/slog"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel/trace"
)

type Observability struct {
	Tracer trace.Tracer
	Logger *slog.Logger
}

func init() {
	ioc.Registry(
		NewObservability,
		newTraceProvider,
		newLoggerProvider)
}
func NewObservability(tracer trace.Tracer, logger *slog.Logger) Observability {
	return Observability{
		Tracer: tracer,
		Logger: logger,
	}
}
