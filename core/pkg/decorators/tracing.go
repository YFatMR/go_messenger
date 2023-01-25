package decorators

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func newSpanName(spanName string) string {
	return "/" + spanName
}

func newDeferFunc(span trace.Span, recordError bool, err error) func() {
	return func() {
		if err != nil && recordError {
			span.RecordError(err)
		}
		span.End()
	}
}

func TraceCallback(ctx context.Context, tracer trace.Tracer,
	spanName string, recordError bool, callback func(context.Context) error) (
	err error,
) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, newSpanName(spanName))
	defer newDeferFunc(span, recordError, err)()
	return callback(ctx)
}

func TraceCallbackWithReturnType[T any](ctx context.Context, tracer trace.Tracer,
	spanName string, recordError bool, callback func(context.Context) (T, error)) (
	_ T, err error,
) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, newSpanName(spanName))
	defer newDeferFunc(span, recordError, err)()
	return callback(ctx)
}

func TraceCallbackWithTwoReturnType[T any, G any](ctx context.Context, tracer trace.Tracer,
	spanName string, recordError bool, callback func(context.Context) (T, G, error)) (
	_ T, _ G, err error,
) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, newSpanName(spanName))
	defer newDeferFunc(span, recordError, err)()
	return callback(ctx)
}
