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
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
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
	restFrontServiceAddress := config.GetStringRequired("REST_SERVICE_ADDRESS")
	grpcFrontServiceAddress := config.GetStringRequired("GRPC_SERVICE_ADDRESS")
	userServiceAddress := config.GetStringRequired("USER_SERVICE_ADDRESS")
	authServiceAddress := config.GetStringRequired("AUTH_SERVICE_ADDRESS")
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

	grpcBackoffConfig := backoff.Config{
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
			Addr:              restFrontServiceAddress,
			Handler:           mux,
		}
		logger.Info(
			"Starting to register REST user server",
			zap.String("REST front server address", restFrontServiceAddress),
			zap.String("gRPC front service address", grpcFrontServiceAddress),
		)

		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if err := proto.RegisterFrontHandlerFromEndpoint(ctx, mux, grpcFrontServiceAddress, opts); err != nil {
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

		listener, err := net.Listen("tcp", grpcFrontServiceAddress)
		if err != nil {
			panic(err)
		}
		logger.Info(
			"Starting to register gRPC user server",
			zap.String("grpc front server address", grpcFrontServiceAddress),
			zap.String("user server address", userServiceAddress),
		)

		logger.Info("Connecting to auth service...", zap.String("address", authServiceAddress))
		authClientOpts := []grpc.DialOption{
			grpc.WithKeepaliveParams(grpcKeepaliveParameters),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
			grpc.WithConnectParams(
				grpc.ConnectParams{
					Backoff: grpcBackoffConfig,
				},
			),
		}
		authServiceClient, err := grpcclients.NewProtobufAuthClient(
			ctx, authServiceAddress, microservicesConnectionTimeout, authClientOpts,
		)
		if err != nil {
			logger.Error("Server stopped! Can't connect to auth service", zap.String("address", authServiceAddress))
			panic(err)
		}

		logger.Info("Connecting to user service...", zap.String("address", userServiceAddress))

		unaryInterceptors := []grpc.UnaryClientInterceptor{
			otelgrpc.UnaryClientInterceptor(),
			interceptors.UnaryAuthInterceptor(
				authServiceClient, grpcAuthorizationHeader, grpcAuthorizationAccountIDHeader,
				grpcAuthorizationUserRoleHeader, logger,
			),
		}

		userClientOpts := []grpc.DialOption{
			grpc.WithKeepaliveParams(grpcKeepaliveParameters),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(
				middleware.ChainUnaryClient(
					unaryInterceptors...,
				),
			),
			grpc.WithConnectParams(
				grpc.ConnectParams{
					Backoff: grpcBackoffConfig,
				},
			),
		}
		userServiceClient, err := grpcclients.NewProtobufUserClient(
			ctx, userServiceAddress, microservicesConnectionTimeout, userClientOpts,
		)
		if err != nil {
			logger.Error("Server stopped! Can't connect to user service", zap.String("address", userServiceAddress))
			panic(err)
		}

		server := frontserver.NewFrontServer(
			authServiceClient, userServiceClient, logger, tracer, grpcBackoffConfig,
		)
		proto.RegisterFrontServer(grpcServer, &server)

		logger.Info("Starting serve gRPC front server")
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()
}
