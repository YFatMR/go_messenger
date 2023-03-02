package ctrace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)

	return span.SpanContext().TraceID().String()
}
