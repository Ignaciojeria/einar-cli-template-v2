package observability

import (
	"archetype/app/shared/configuration"
	"archetype/app/shared/infrastructure/observability/strategy"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	ioc.Registry(
		newTraceProvider,
		configuration.NewEnvLoader,
	)
}

// RegisterTraceProvider determines whether to use OpenObserve, Datadog or non provider based on the existing environment variables.
func newTraceProvider(env configuration.EnvLoader) (trace.Tracer, error) {
	// Get the observability strategy
	observabilityStrategyKey := env.Get("OBSERVABILITY_STRATEGY")
	switch observabilityStrategyKey {
	case "openobserve":
		return strategy.OpenObserveGRPCTraceProvider(env)
	case "datadog":
		return strategy.DatadogGRPCTraceProvider(env)
	default:
		return strategy.NoOpTraceProvider(env)
	}
}
