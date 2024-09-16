package strategy

import (
	"archetype/app/shared/configuration"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// newGRPCOpenObserveLoggerProvider configures the logger provider for OpenObserve.
func OpenObserveGRPCLogProvider(env configuration.EnvLoader) (*slog.Logger, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Configure the exporter options
	exporterOpts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(env.Get("OTEL_EXPORTER_OTLP_ENDPOINT")),
		otlploggrpc.WithTLSCredentials(insecure.NewCredentials()),
		otlploggrpc.WithDialOption(grpc.WithUnaryInterceptor(func(
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
	}

	// Create the exporter
	exporter, err := otlploggrpc.New(ctx, exporterOpts...)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("creating OTLP log exporter: %w", err)
	}

	// Set up the processor
	logProcessor := log.NewBatchProcessor(exporter)

	// Define the resource attributes
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(env.Get("PROJECT_NAME")),
		semconv.ServiceVersionKey.String(env.Get("VERSION")),
		semconv.DeploymentEnvironmentKey.String(env.Get("ENVIRONMENT")),
	)

	// Create the LoggerProvider
	loggerProvider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(logProcessor),
	)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCancel()
		if err := loggerProvider.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown logger provider:", err)
		}
		cancel()
	}()

	// Create the slog.Logger using the otelslog bridge
	logger := otelslog.NewLogger(
		"openobserve",
		otelslog.WithLoggerProvider(loggerProvider),
	)

	return logger, nil
}
