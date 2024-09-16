package strategy

import (
	"archetype/app/shared/configuration"
	"context"
	"log/slog"
	"os"
	"strconv"

	"go.opentelemetry.io/otel/trace"
)

const (
	ddTraceIDKey = "dd.trace_id"
	ddSpanIDKey  = "dd.span_id"
	ddServiceKey = "dd.service"
	ddEnvKey     = "dd.env"
	ddVersionKey = "dd.version"
)

func DatadogStdoutLogProvider(env configuration.EnvLoader) *slog.Logger {
	// Crear el baseHandler que imprima los logs en formato JSON en consola
	baseHandler := slog.NewJSONHandler(os.Stdout, nil)

	// Crear el DatadogHandler para capturar datos del contexto y añadirlos a los logs
	datadogHandler := newDatadogHandler(baseHandler)

	// Crear el logger utilizando el handler personalizado
	return slog.New(datadogHandler).With(
		slog.String(ddEnvKey, env.Get("ENVIRONMENT")),
		slog.String(ddVersionKey, env.Get("VERSION")),
		slog.String(ddServiceKey, env.Get("PROJECT_NAME")),
	)
}

// DatadogHandler es un handler personalizado que extrae datos adicionales del contexto y los envía a Datadog.
type DatadogHandler struct {
	baseHandler slog.Handler
}

func newDatadogHandler(baseHandler slog.Handler) *DatadogHandler {
	return &DatadogHandler{baseHandler: baseHandler}
}

// Handle implementa el método para manejar los logs, añadiendo elementos del contexto.
func (h *DatadogHandler) Handle(ctx context.Context, record slog.Record) error {
	// Extraer traceID y spanID desde el contexto (OpenTelemetry)
	if spanContext := trace.SpanContextFromContext(ctx); spanContext.IsValid() {
		record.AddAttrs(
			slog.String(ddTraceIDKey, convertTraceID(spanContext.TraceID().String())),
			slog.String(ddSpanIDKey, convertTraceID(spanContext.SpanID().String())),
		)
	}

	// Delegar al baseHandler para que continúe el procesamiento del log
	return h.baseHandler.Handle(ctx, record)
}

// WithAttrs permite añadir atributos adicionales al handler.
func (h *DatadogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &DatadogHandler{baseHandler: h.baseHandler.WithAttrs(attrs)}
}

// WithGroup permite añadir grupos de atributos al handler.
func (h *DatadogHandler) WithGroup(name string) slog.Handler {
	return &DatadogHandler{baseHandler: h.baseHandler.WithGroup(name)}
}

// Enabled verifica si el log debe ser manejado basado en el nivel de logs.
func (h *DatadogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Delegar al baseHandler
	return h.baseHandler.Enabled(ctx, level)
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
