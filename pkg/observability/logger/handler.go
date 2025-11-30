package logger

import (
	"context"
	"log/slog"
)

type NoopHandler struct{}

func (h *NoopHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

func (h *NoopHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func (h *NoopHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *NoopHandler) WithGroup(name string) slog.Handler {
	return h
}
