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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// init function to conditionally register the appropriate trace provider
func init() {
	ioc.Registry(
		NewTraceProvider,
		configuration.NewEnvLoader,
	)
}

// RegisterTraceProvider determines whether to use OpenObserve or a basic provider based on the existing environment variables.
func NewTraceProvider(env configuration.EnvLoader) (trace.Tracer, error) {
	// Check if OpenObserve variables are set
	if isOpenObserveConfigured(env) {
		return newGRPCOpenObserveTraceProvider(env)
	}
	return newDefaultTraceProvider(env)
}

// isOpenObserveConfigured checks if the environment variables for OpenObserve are set.
func isOpenObserveConfigured(env configuration.EnvLoader) bool {
	endpoint := env.Get("OPENOBSERVE_GRPC_ENDPOINT")
	auth := env.Get("OPENOBSERVE_AUTHORIZATION")
	org := env.Get("OPENOBSERVE_ORGANIZATION")
	stream := env.Get("OPENOBSERVE_STREAM_NAME")
	// If all the key variables for OpenObserve are present, we use OpenObserve
	return endpoint != "" && auth != "" && org != "" && stream != ""
}

// NewGRPCOpenObserveTraceProvider configures the trace provider for OpenObserve.
func newGRPCOpenObserveTraceProvider(env configuration.EnvLoader) (trace.Tracer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(env.Get("OPENOBSERVE_GRPC_ENDPOINT")),
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

// NewTraceProvider configures a basic trace provider without OpenObserve.
func newDefaultTraceProvider(env configuration.EnvLoader) (trace.Tracer, error) {
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
