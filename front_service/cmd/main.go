package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	grpcclients "github.com/YFatMR/go_messenger/core/pkg/grpc_clients"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/traces"
	"github.com/YFatMR/go_messenger/front_server/internal/interceptors"
	frontserver "github.com/YFatMR/go_messenger/front_server/internal/server"
	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
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
	authServerAddress := config.GetStringRequired("AUTH_SERVICE_ADDRESS")
	jaegerEndpoint := config.GetStringRequired("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := config.GetStringRequired("SERVICE_NAME")
	restServiceReadTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_READ_TIMEOUT_SECONDS")
	restServiceWriteTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_WRITE_TIMEOUT_SECONDS")
	restServiceIdleTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_IDLE_TIMEOUT_SECONDS")
	restServiceReadHeaderTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_READ_HEADER_TIMEOUT_SECONDS")
	microservicesConnectionTimeout := config.GetMillisecondsDurationRequired(
		"MICROSERVICES_GRPC_CONNECTION_TIMEOUT_MILLISECONDS",
	)

	grpcAuthorizationHeader := config.GetStringRequired("GRPC_AUTHORIZARION_HEADER")
	grpcAuthorizationAccountIDHeader := config.GetStringRequired("GRPC_AUTHORIZARION_ACCOUNT_ID_HEADER")
	grpcAuthorizationUserRoleHeader := config.GetStringRequired("GRPC_AUTHORIZARION_USER_ROLE_HEADER")

	grpcBackoffCongig := backoff.Config{
		BaseDelay:  config.GetMillisecondsDurationRequired("GRPC_CONNECTION_BACKOFF_DELAY_MILLISECONDS"),
		Multiplier: config.GetFloat64Required("GRPC_CONNECTION_BACKOFF_MULTIPLIER"),
		Jitter:     config.GetFloat64Required("GRPC_CONNECTION_BACKOFF_JITTER"),
		MaxDelay:   config.GetMillisecondsDurationRequired("GRPC_CONNECTION_BACKOFF_MAX_DELAY_MILLISECONDS"),
	}

	grpcKeepaliveParameters := keepalive.ClientParameters{
		Time:                config.GetMillisecondsDurationRequired("GRPC_CONNECTION_KEEPALIVE_TIME_MILLISECONDS"),
		Timeout:             config.GetMillisecondsDurationRequired("GRPC_CONNECTION_KEEPALIVE_TIMEOUT_MILLISECONDS"),
		PermitWithoutStream: config.GetBoolRequired("GRPC_CONNECTION_KEEPALIVE_PERMIT_WITHOUT_STREAM"),
	}

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

		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if err := proto.RegisterFrontUserHandlerFromEndpoint(ctx, mux, grpcFrontUserServerAddress, opts); err != nil {
			panic(err)
		}

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

		logger.Info("Connecting to auth service...", zap.String("address", authServerAddress))
		authServiceClient, err := grpcclients.NewProtobufAuthClient(
			ctx, authServerAddress, microservicesConnectionTimeout, grpcBackoffCongig, grpcKeepaliveParameters,
		)
		if err != nil {
			logger.Error("Server stopped! Can't connect to auth service", zap.String("address", authServerAddress))
			panic(err)
		}

		logger.Info("Connecting to user service...", zap.String("address", userServerAddress))

		unaryInterceptors := []grpc.UnaryClientInterceptor{
			otelgrpc.UnaryClientInterceptor(),
			interceptors.UnaryAuthInterceptor(
				authServiceClient, grpcAuthorizationHeader, grpcAuthorizationAccountIDHeader,
				grpcAuthorizationUserRoleHeader, logger,
			),
		}

		userServiceClient, err := grpcclients.NewProtobufUserClient(
			ctx, userServerAddress, microservicesConnectionTimeout, grpcBackoffCongig, grpcKeepaliveParameters,
			unaryInterceptors,
		)
		if err != nil {
			logger.Error("Server stopped! Can't connect to user service", zap.String("address", userServerAddress))
			panic(err)
		}

		server := frontserver.NewFrontUserServer(
			authServiceClient, userServiceClient, logger, tracer, grpcBackoffCongig,
		)
		proto.RegisterFrontUserServer(grpcServer, &server)

		logger.Info("Starting serve gRPC front server")
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()
}
