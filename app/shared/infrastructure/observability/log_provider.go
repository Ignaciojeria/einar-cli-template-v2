package observability

import (
	"archetype/app/shared/configuration"
	"archetype/app/shared/infrastructure/observability/strategy"
	"log/slog"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		newLoggerProvider,
		configuration.NewEnvLoader,
	)
}

func newLoggerProvider(env configuration.EnvLoader) (*slog.Logger, error) {
	// Get the observability strategy
	observabilityStrategyKey := env.Get("OBSERVABILITY_STRATEGY")
	switch observabilityStrategyKey {
	case "openobserve":
		return strategy.OpenObserveGRPCLogProvider(env)
	case "datadog":
		return strategy.DatadogStdoutLogProvider(env), nil
	default:
		return strategy.NoOpStdoutLogProvider(env), nil
	}
}
