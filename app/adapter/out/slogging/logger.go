package slogging

import (
	"log/slog"
	"os"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type Logger struct {
	*slog.Logger
}

func init() {
	ioc.Registry(NewLogger)
}
func NewLogger() Logger {
	return Logger{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}
