package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/traces"
	userserver "github.com/YFatMR/go_messenger/front_server/internal/user_server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := cviper.New()
	config.AutomaticEnv()

	// Init environment vars
	logLevel := config.GetZapcoreLogLevelRequired("LOG_LEVEL")
	logPath := config.GetStringRequired("LOG_PATH")
	restFrontUserServerAddress := config.GetStringRequired("REST_SERVICE_ADDRESS")
	grpcFrontUserServerAddress := config.GetStringRequired("GRPC_SERVICE_ADDRESS")
	userServerAddress := config.GetStringRequired("USER_SERVICE_ADDRESS")
	jaegerEndpoint := config.GetStringRequired("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := config.GetStringRequired("SERVICE_NAME")
	restServiceReadTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_READ_TIMEOUT_SECONDS")
	restServiceWriteTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_WRITE_TIMEOUT_SECONDS")
	restServiceIdleTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_IDLE_TIMEOUT_SECONDS")
	restServiceReadHeaderTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_READ_HEADER_TIMEOUT_SECONDS")

	// Init logger
	zapLogger, err := loggers.NewBaseZapFileLogger(logLevel, logPath)
	if err != nil {
		panic(err)
	}
	logger := loggers.NewOtelZapLoggerWithTraceID(
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
	defer func(ctx context.Context) {
		if err := traceProvider.Shutdown(ctx); err != nil { // TODO: Shutdown -> with timeout
			panic(err)
		}
	}(ctx)
	tracer := otel.Tracer(serviceName)

	// Run REST front service with "user service" API
	go func() {
		mux := runtime.NewServeMux()

		service := http.Server{
			ReadTimeout:       restServiceReadTimeout,
			WriteTimeout:      restServiceWriteTimeout,
			IdleTimeout:       restServiceIdleTimeout,
			ReadHeaderTimeout: restServiceReadHeaderTimeout,
			Addr:              restFrontUserServerAddress,
			Handler:           mux,
		}
		logger.Info(
			"Starting to register REST user server",
			zap.String("REST front server address", restFrontUserServerAddress),
			zap.String("gRPC front service address", grpcFrontUserServerAddress),
		)
		userserver.RegisterRestUserServer(ctx, mux, grpcFrontUserServerAddress)
		logger.Info("Starting serve REST front server")
		//#nosec G114: Use of net/http serve function that has no support for setting timeouts
		if err := service.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	// Run gRPC front service with "user service" API
	func() {
		grpcServer := grpc.NewServer()

		listener, err := net.Listen("tcp", grpcFrontUserServerAddress)
		if err != nil {
			panic(err)
		}
		logger.Info(
			"Starting to register gRPC user server",
			zap.String("grpc front server address", grpcFrontUserServerAddress),
			zap.String("user server address", userServerAddress),
		)
		userserver.RegisterGrpcUserServer(grpcServer, userServerAddress, logger, tracer)
		logger.Info("Starting serve gRPC front server")
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()
}
