package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/traces"
	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/controllers"
	mongorepository "github.com/YFatMR/go_messenger/user_service/internal/repositories/mongo"
	"github.com/YFatMR/go_messenger/user_service/internal/servers"
	"github.com/YFatMR/go_messenger/user_service/internal/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func main() {
	time.Sleep(3 * time.Second)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := cviper.New()
	config.AutomaticEnv()

	// Init environment vars
	logLevel := config.GetZapcoreLogLevelRequired("LOG_LEVEL")
	logPath := config.GetStringRequired("LOG_PATH")
	mongoURI := config.GetStringRequired("MONGODB_URI")
	databaseName := config.GetStringRequired("MONGODB_DATABASE_NAME")
	collectionName := config.GetStringRequired("MONGODB_DATABASE_COLLECTION_NAME")
	mongoOperationTimeout := config.GetSecondsDurationRequired("MONGODB_OPERATION_TIMEOUT_SECONDS")
	jaegerEndpoint := config.GetStringRequired("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := config.GetStringRequired("SERVICE_NAME")
	userServiceAddress := config.GetStringRequired("SERVICE_ADDRESS")
	metricsServiceReadTimeout := config.GetSecondsDurationRequired("METRICS_SERVICE_READ_TIMEOUT_SECONDS")
	metricsServiceWriteTimeout := config.GetSecondsDurationRequired("METRICS_SERVICE_WRITE_TIMEOUT_SECONDS")
	metricsServiceIdleTimeout := config.GetSecondsDurationRequired("METRICS_SERVICE_IDLE_TIMEOUT_SECONDS")
	metricsServiceReadHeaderTimeout := config.GetSecondsDurationRequired("METRICS_SERVICE_READ_HEADER_TIMEOUT_SECONDS")
	metricsServiceListingSuffix := config.GetStringRequired("METRICS_SERVICE_LISTING_SUFFIX")
	metricsServiceAddress := config.GetStringRequired("METRICS_SERVICE_ADDRESS")

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

	// Init database
	logger.Info("Init database")
	mongoCollection, cancelConnection := func() (*mongo.Collection, func()) {
		ctx, cancel := context.WithTimeout(ctx, mongoOperationTimeout)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			panic(err)
		}

		logger.Info("Starting Ping mongodb")
		err = client.Ping(ctx, nil)
		if err != nil {
			panic(err)
		}
		logger.Info("mongodb Ping successfully finished")

		collection := client.Database(databaseName).Collection(collectionName)
		return collection, cancel
	}()
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

	go func(logger *loggers.OtelZapLoggerWithTraceID) {
		server := &http.Server{
			ReadTimeout:       metricsServiceReadTimeout,
			WriteTimeout:      metricsServiceWriteTimeout,
			IdleTimeout:       metricsServiceIdleTimeout,
			ReadHeaderTimeout: metricsServiceReadHeaderTimeout,
			Addr:              metricsServiceAddress,
			Handler:           nil,
		}
		http.Handle(metricsServiceListingSuffix, promhttp.Handler())
		//#nosec G114: Use of net/http serve function that has no support for setting timeouts
		if err := server.ListenAndServe(); err != nil {
			logger.Error("Can't up metrics server with endpoint" + metricsServiceAddress +
				". Operation finished with error: " + err.Error())
			panic(err)
		}
	}(logger)

	// Init Server
	logger.Info("Init server")
	userRepository := mongorepository.NewUserMongoRepository(mongoCollection, mongoOperationTimeout, logger, tracer)
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
