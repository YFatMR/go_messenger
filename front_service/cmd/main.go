package main

import (
	"context"
	. "github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/traces"
	. "github.com/YFatMR/go_messenger/core/pkg/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func main() {

	// Init environment vars
	logLevel := RequiredZapcoreLogLevelEnv("LOG_LEVEL")
	logPath := RequiredStringEnv("LOG_PATH")
	frontRestUserServerAddress := RequiredStringEnv("REST_SERVICE_ADDRESS")
	frontGrpcUserServerAddress := RequiredStringEnv("GRPC_SERVICE_ADDRESS")
	userServerAddress := RequiredStringEnv("USER_SERVICE_ADDRESS")
	jaegerEndpoint := RequiredStringEnv("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := RequiredStringEnv("SERVICE_NAME")

	// Init logger
	zapLogger, err := NewBaseZapFileLogger(logLevel, logPath)
	if err != nil {
		panic(err)
	}
	logger := NewOtelZapLoggerWithTraceID(
		otelzap.New(
			zapLogger,
			otelzap.WithTraceIDField(true),
			otelzap.WithMinLevel(zapcore.ErrorLevel),
			otelzap.WithStackTrace(true),
		),
	)
	defer logger.Sync()

	// Init tracing
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
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
