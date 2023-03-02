package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/ctrace"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
	"github.com/YFatMR/go_messenger/core/pkg/metrics/prometheus"
	"github.com/YFatMR/go_messenger/core/pkg/mongodb"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/controllers"
	"github.com/YFatMR/go_messenger/user_service/controllers/controllerdecorators"
	"github.com/YFatMR/go_messenger/user_service/grpcserver"
	"github.com/YFatMR/go_messenger/user_service/mongorepository"
	"github.com/YFatMR/go_messenger/user_service/passwordmanager"
	"github.com/YFatMR/go_messenger/user_service/repositories"
	"github.com/YFatMR/go_messenger/user_service/repositories/repositorydecorators"
	"github.com/YFatMR/go_messenger/user_service/services"
	"github.com/YFatMR/go_messenger/user_service/services/servicedecorators"
	"github.com/YFatMR/go_messenger/user_service/usercontroller"
	"github.com/YFatMR/go_messenger/user_service/userservice"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := cviper.New()
	config.AutomaticEnv()

	// Init environment vars
	databaseSettings := cviper.NewDatabaseSettingsFromConfig(config)
	metricServiceSettings := cviper.NewMetricMetricServiceSettingsFromConfig(config)

	jaegerEndpoint := config.GetStringRequired("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := config.GetStringRequired("SERVICE_NAME")
	serviceAddress := config.GetStringRequired("SERVICE_ADDRESS")

	collectDatabaseQueryMetrics := config.GetBoolRequired("ENABLE_DATABASE_QUERY_METRICS")
	traceDatabaseQuery := config.GetBoolRequired("ENABLE_DATABASE_QUERY_TRACING")

	// Init logger
	logger, err := czap.FromConfig(config)
	if err != nil {
		panic(err)
	}
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
	traceProvider, err := ctrace.NewProvider(exporter, resourcesdk.NewWithAttributes(
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
	repository = repositorydecorators.NewLoggingUserRepositoryDecorator(repository, logger)
	if collectDatabaseQueryMetrics {
		repository = repositorydecorators.NewPrometheusMetricsUserRepositoryDecorator(repository)
	}
	if traceDatabaseQuery {
		// TODO: make as config option (?)
		recordTraceErrors := true
		repository = repositorydecorators.NewOpentelemetryTracingUserRepositoryDecorator(
			repository, tracer, recordTraceErrors,
		)
	}

	passwordManager := passwordmanager.Default()
	jwtManager := jwtmanager.FromConfig(config, logger)
	var service services.UserService
	service = userservice.New(repository, passwordManager, jwtManager, logger)
	// TODO: make config options
	if true {
		service = servicedecorators.NewLoggingUserServiceDecorator(service, logger)
	}
	if true {
		service = servicedecorators.NewPrometheusMetricsUserServiceDecorator(service)
	}
	if true {
		recordTraceErrors := true
		service = servicedecorators.NewOpentelemetryTracingUserServiceDecorator(service, tracer, recordTraceErrors)
	}

	var controller controllers.UserController
	controller = usercontroller.New(service, logger)
	// TODO: make config options
	if true {
		controller = controllerdecorators.NewLoggingUserControllerDecorator(controller, logger)
	}

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
