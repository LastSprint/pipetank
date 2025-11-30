package ctx

import (
	"context"
	"log/slog"
)

type loggerKeyType string

const (
	loggerKey loggerKeyType = "logger"
)

func Logger(ctx context.Context) *slog.Logger {
	vl := ctx.Value(loggerKey)

	if vl == nil {
		return slog.Default()
	}

	w, ok := vl.(*slog.Logger)
	if !ok {
		return slog.Default()
	}

	return w
}

func SetLogger(ctx context.Context, w *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, w)
}
