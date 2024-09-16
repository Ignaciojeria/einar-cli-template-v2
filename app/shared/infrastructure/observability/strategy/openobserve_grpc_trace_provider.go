package strategy

import (
	"archetype/app/shared/configuration"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// NewGRPCOpenObserveTraceProvider configures the trace provider for OpenObserve.
func OpenObserveGRPCTraceProvider(env configuration.EnvLoader) (trace.Tracer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(env.Get("OTEL_EXPORTER_OTLP_ENDPOINT")),
		otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()),
		otlptracegrpc.WithDialOption(grpc.WithUnaryInterceptor(func(
			ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
			opts ...grpc.CallOption) error {
			md := metadata.New(map[string]string{
				"Authorization": env.Get("OPENOBSERVE_AUTHORIZATION"),
				"organization":  env.Get("OPENOBSERVE_ORGANIZATION"),
				"stream-name":   env.Get("OPENOBSERVE_STREAM_NAME"),
			})
			ctx = metadata.NewOutgoingContext(ctx, md)
			return invoker(ctx, method, req, reply, cc, opts...)
		})),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(env.Get("PROJECT_NAME")),
			semconv.DeploymentEnvironmentKey.String(env.Get("ENVIRONMENT")),
		)),
	)

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
