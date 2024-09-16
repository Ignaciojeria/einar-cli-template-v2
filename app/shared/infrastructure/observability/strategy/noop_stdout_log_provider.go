package strategy

import (
	"archetype/app/shared/configuration"
	"log/slog"
	"os"

	otelslogjson "github.com/go-slog/otelslog"
)

func NoOpStdoutLogProvider(env configuration.EnvLoader) *slog.Logger {
	return slog.New(otelslogjson.NewHandler(slog.NewJSONHandler(os.Stdout, nil))).With(
		slog.String("env", env.Get("ENVIRONMENT")),
		slog.String("version", env.Get("VERSION")),
		slog.String("service", env.Get("PROJECT_NAME")),
	)
}
