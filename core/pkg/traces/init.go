package traces

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func NewTraceProvider(exporter tracesdk.SpanExporter, resource *resourcesdk.Resource) (*tracesdk.TracerProvider, error) {
	//exporter, err := newJaegerExporter("http://" + jaegerEndpoint + "/api/traces")
	//if err != nil {
	//	return nil, err
	//}
	// resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("Front service"))
	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return traceProvider, nil
}
