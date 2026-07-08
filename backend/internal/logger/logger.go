package logger

import (
	"log/slog"
	"os"
)

// New creates a structured logger with the given environment and base attributes.
// In development mode it uses human-readable text output; otherwise JSON.
func New(env string, attrs ...slog.Attr) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if env == "development" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	l := slog.New(handler)
	for _, a := range attrs {
		l = l.With(a.Key, a.Value.Any())
	}
	return l
}
