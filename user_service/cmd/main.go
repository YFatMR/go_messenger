package main

import (
	"context"
	. "core/pkg/loggers"
	"core/pkg/traces"
	. "core/pkg/utils"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"net"
	proto "protocol/pkg/proto"
	"time"
	"user_server/internal/controllers"
	"user_server/internal/repositories/mongo"
	"user_server/internal/servers"
	"user_server/internal/services"
)

// GOMAXPROC

func newJaegerExporter(jaegerEndpoint string) (tracesdk.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
}

func main() {
	time.Sleep(5 * time.Second)

	// Init environment vars
	logLevel := RequiredZapcoreLogLevelEnv("USER_SERVICE_LOG_LEVEL")
	logPath := RequiredStringEnv("USER_SERVICE_LOG_PATH")
	mongoUri := RequiredStringEnv("USER_SERVICE_MONGODB_URI")
	databaseName := RequiredStringEnv("USER_SERVICE_MONGODB_DATABASE_NAME")
	collectionName := RequiredStringEnv("USER_SERVICE_MONGODB_DATABASE_COLLECTION_NAME")
	connectionTimeout := RequiredIntEnv("USER_SERVICE_MONGODB_CONNECTION_TIMEOUT_SEC")
	jaegerEndpoint := RequiredStringEnv("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := RequiredStringEnv("SERVICE_NAME")
	userServerAddress := GetFullServiceAddress("USER")
	//fmt.Println("some testing,,,,", viper.GetString("USER_SERVICE_LOG_PATH"))

	// Init logger
	logger, err := NewOtelZapLoggerWithTraceID(logLevel, logPath)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Init database
	logger.Info("Init database")
	mongoSetting := mongo.NewMongoSettings(mongoUri, databaseName, collectionName, time.Duration(connectionTimeout)*time.Second)
	mongoCtx := context.Background()
	mongoCollection, cancelConnection := mongo.NewMongoCollection(mongoCtx, mongoSetting, logger)
	defer cancelConnection()

	// Init metrics
	logger.Info("Init metrics")
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

	// Init Server
	logger.Info("Init server")
	userRepository := mongo.NewUserMongoRepository(mongoCollection, logger)
	userService := services.NewUserService(userRepository, logger)
	userController := controllers.NewUserController(userService, logger)
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
	listener, err := net.Listen("tcp", userServerAddress)
	if err != nil {
		panic(err)
	}
	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
