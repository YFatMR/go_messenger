package main

import (
	"context"
	. "core/pkg/loggers"
	"core/pkg/traces"
	. "core/pkg/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
)

func main() {

	// Init environment vars
	logLevel := RequiredZapcoreLogLevelEnv("FRONT_SERVICE_LOG_LEVEL")
	logPath := RequiredStringEnv("FRONT_SERVICE_LOG_PATH")
	frontRestUserServerAddress := GetFullServiceAddress("FRONT_REST")
	frontGrpcUserServerAddress := GetFullServiceAddress("FRONT_GRPC")
	userServerAddress := GetFullServiceAddress("USER")
	jaegerEndpoint := RequiredStringEnv("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := RequiredStringEnv("SERVICE_NAME")

	// Init logger
	logger, err := NewOtelZapLoggerWithTraceID(logLevel, logPath)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Init metrics
	exporter, err := traces.NewJaegerExporter(jaegerEndpoint)
	if err != nil {
		panic(err)
	}
	traceProvider, err := traces.NewTraceProvider(exporter, resourcesdk.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil { // TODO: Shutdown -> with timeout
			panic(err)
		}
	}()
	tracer := otel.Tracer(serviceName)

	// Init servers
	mux := runtime.NewServeMux()
	grpcServer := grpc.NewServer()
	ctx := context.Background()

	// Run servers
	go runRest(ctx, mux, frontRestUserServerAddress, frontGrpcUserServerAddress, logger)
	runGrpc(grpcServer, frontGrpcUserServerAddress, userServerAddress, logger, tracer)
}
