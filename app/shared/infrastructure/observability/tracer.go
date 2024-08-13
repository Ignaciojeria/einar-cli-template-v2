package observability

import (
	"archetype/app/shared/configuration"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	ioc.Registry(
		NewTraceProvider,
		configuration.NewConf)
}

func NewTraceProvider(conf configuration.Conf) (trace.Tracer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	client := otlptracegrpc.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.PROJECT_NAME),
			semconv.DeploymentEnvironmentKey.String(conf.ENVIRONMENT),
		)),
	)
	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Second*2)
		defer shutdownCancel()
		if err := tp.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown:", err)
		}
		cancel()
	}()
	return tp.Tracer("observability"), nil
}
