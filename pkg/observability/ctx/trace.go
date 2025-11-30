package ctx

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func GetTraceID(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).TraceID().String()
}
