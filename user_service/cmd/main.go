package main

import (
	"context"
	. "github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/traces"
	. "github.com/YFatMR/go_messenger/core/pkg/utils"
	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/controllers"
	"github.com/YFatMR/go_messenger/user_service/internal/repositories/mongo"
	"github.com/YFatMR/go_messenger/user_service/internal/servers"
	"github.com/YFatMR/go_messenger/user_service/internal/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"
)

func main() {
	time.Sleep(5 * time.Second)

	// Init environment vars
	logLevel := RequiredZapcoreLogLevelEnv("LOG_LEVEL")
	logPath := RequiredStringEnv("LOG_PATH")
	mongoUri := RequiredStringEnv("MONGODB_URI")
	databaseName := RequiredStringEnv("MONGODB_DATABASE_NAME")
	collectionName := RequiredStringEnv("MONGODB_DATABASE_COLLECTION_NAME")
	connectionTimeout := RequiredIntEnv("MONGODB_CONNECTION_TIMEOUT_SEC")
	jaegerEndpoint := RequiredStringEnv("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := RequiredStringEnv("SERVICE_NAME")
	userServiceAddress := RequiredStringEnv("SERVICE_ADDRESS")

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

	// Init database
	logger.Info("Init database")
	mongoSetting := mongo.NewMongoSettings(mongoUri, databaseName, collectionName, time.Duration(connectionTimeout)*time.Second)
	mongoCtx := context.Background()
	mongoCollection, cancelConnection := mongo.NewMongoCollection(mongoCtx, mongoSetting, logger)
	defer cancelConnection()

	// Init tracing
	logger.Info("Init metrics")
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

	// Init metrics

	go func(logger *OtelZapLoggerWithTraceID) {
		http.Handle(RequiredStringEnv("METRICS_LISTING_SUFFIX"), promhttp.Handler())
		metricsEndpoint := RequiredStringEnv("METRICS_ADDRESS")
		if err := http.ListenAndServe(metricsEndpoint, nil); err != nil {
			logger.Error("Can't up metrics server with endpoint" + metricsEndpoint + ". Operation finished with error: " + err.Error())
			panic(err)
		}
	}(logger)

	// Init Server
	logger.Info("Init server")
	userRepository := mongo.NewUserMongoRepository(mongoCollection, logger, tracer)
	userService := services.NewUserService(userRepository, logger, tracer)
	userController := controllers.NewUserController(userService, logger, tracer)
	gRPCUserServer := servers.NewGRPCUserServer(userController, logger, tracer)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	// Register protobuf server
	logger.Info("Register protobuf server")
	proto.RegisterUserServer(s, gRPCUserServer)

	// Listen connection
	logger.Info("Server successfully setup. Starting listen...")
	listener, err := net.Listen("tcp", userServiceAddress)
	if err != nil {
		panic(err)
	}
	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
