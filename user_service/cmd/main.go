package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/metrics/prometheus"
	"github.com/YFatMR/go_messenger/core/pkg/mongodb"
	"github.com/YFatMR/go_messenger/core/pkg/traces"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/controllers"
	cdecorators "github.com/YFatMR/go_messenger/user_service/internal/controllers/decorators"
	"github.com/YFatMR/go_messenger/user_service/internal/controllers/usercontroller"
	"github.com/YFatMR/go_messenger/user_service/internal/repositories"
	rdecorators "github.com/YFatMR/go_messenger/user_service/internal/repositories/decorators"
	"github.com/YFatMR/go_messenger/user_service/internal/repositories/mongorepository"
	"github.com/YFatMR/go_messenger/user_service/internal/servers/grpcserver"
	"github.com/YFatMR/go_messenger/user_service/internal/services"
	sdecorators "github.com/YFatMR/go_messenger/user_service/internal/services/decorators"
	"github.com/YFatMR/go_messenger/user_service/internal/services/userservice"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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

	databaseSettings := cviper.NewDatabaseSettingsFromConfig(config)
	metricServiceSettings := cviper.NewMetricMetricServiceSettingsFromConfig(config)

	jaegerEndpoint := config.GetStringRequired("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := config.GetStringRequired("SERVICE_NAME")
	serviceAddress := config.GetStringRequired("SERVICE_ADDRESS")

	collectDatabaseQueryMetrics := config.GetBoolRequired("ENABLE_DATABASE_QUERY_METRICS")
	traceDatabaseQuery := config.GetBoolRequired("ENABLE_DATABASE_QUERY_TRACING")

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
	mongoCollection, err := mongodb.Connect(ctx, databaseSettings, logger)
	if err != nil {
		logger.Error("Can't establish connection with database", zap.Error(err))
		panic(err)
	}

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
	if metricServiceSettings != nil {
		go prometheus.ListenAndServeMetrcirService(metricServiceSettings, logger)
	}

	// Init Server
	logger.Info("Init server")

	var repository repositories.UserRepository
	repository = mongorepository.NewUserMongoRepository(mongoCollection, databaseSettings.GetOperationTimeout(), logger)
	repository = rdecorators.NewLoggingUserRepositoryDecorator(repository, logger)
	if collectDatabaseQueryMetrics {
		repository = rdecorators.NewPrometheusMetricsUserRepositoryDecorator(repository)
	}
	if traceDatabaseQuery {
		// TODO: make as config option (?)
		recordTraceErrors := true
		repository = rdecorators.NewOpentelemetryTracingUserRepositoryDecorator(repository, tracer, recordTraceErrors)
	}
	repository = rdecorators.NewLogerrCleanerUserRepositoryDecorator(repository)

	var service services.UserService
	service = userservice.New(repository)
	// TODO: make config options
	if true {
		service = sdecorators.NewLoggingUserServiceDecorator(service, logger)
	}
	if true {
		repository = sdecorators.NewPrometheusMetricsUserServiceDecorator(service)
	}
	if true {
		recordTraceErrors := true
		service = sdecorators.NewOpentelemetryTracingUserServiceDecorator(service, tracer, recordTraceErrors)
	}
	service = sdecorators.NewLogerrCleanerUserServiceDecorator(service)

	var controller controllers.UserController
	controller = usercontroller.New(service)
	// TODO: make config options
	if true {
		controller = cdecorators.NewLoggingUserControllerDecorator(controller, logger)
	}
	controller = cdecorators.NewLogerrCleanerUserControllerDecorator(controller)

	server := grpcserver.New(controller)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	// Register protobuf server
	logger.Info("Register protobuf server")
	proto.RegisterUserServer(s, &server)

	// Listen connection
	logger.Info("Server successfully setup. Starting listen...")
	listener, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		panic(err)
	}
	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
