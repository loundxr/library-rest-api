package logger

import (
	"context"
	"fmt"
	"library-rest-api/internal/logger/slogpretty"
	"log/slog"
	"os"
)

type TeeHandler struct {
	handlers []slog.Handler
}

func (t *TeeHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range t.handlers {
		if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func (t *TeeHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return t.handlers[0].Enabled(ctx, l)
}

func (t *TeeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(t.handlers))
	for i, h := range t.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &TeeHandler{handlers: newHandlers}
}

func (t *TeeHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(t.handlers))
	for i, h := range t.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &TeeHandler{handlers: newHandlers}
}

func New() (*slog.Logger, *os.File, error) {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %w", err)
	}

	fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	prettyOpts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	consoleHandler := prettyOpts.NewPrettyHandler(os.Stdout)

	combinedHandler := &TeeHandler{
		handlers: []slog.Handler{fileHandler, consoleHandler},
	}

	logger := slog.New(combinedHandler)
	return logger, logFile, nil
}
