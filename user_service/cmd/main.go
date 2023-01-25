package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"github.com/YFatMR/go_messenger/core/pkg/mongodb"
	"github.com/YFatMR/go_messenger/core/pkg/traces"
	proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/controllers"
	"github.com/YFatMR/go_messenger/user_service/internal/repositories"
	"github.com/YFatMR/go_messenger/user_service/internal/repositories/decorators"
	"github.com/YFatMR/go_messenger/user_service/internal/repositories/mongorepository"
	"github.com/YFatMR/go_messenger/user_service/internal/servers"
	"github.com/YFatMR/go_messenger/user_service/internal/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	databaseURI := config.GetStringRequired("DATABASE_URI")
	databaseName := config.GetStringRequired("DATABASE_NAME")
	databaseCollectionName := config.GetStringRequired("DATABASE_COLLECTION_NAME")
	databaseOperationTimeout := config.GetSecondsDurationRequired("DATABASE_OPERATION_TIMEOUT_MILLISECONDS")
	databaseConnectionTimeout := config.GetMillisecondsDurationRequired("DATABASE_CONNECTION_TIMEOUT_MILLISECONDS")
	databaseReconnectionCount := config.GetIntRequired("DATABASE_RECONNECTION_COUNT")
	databaseReconnectionInterval := config.GetMillisecondsDurationRequired(
		"DATABASE_RECONNECTIONION_INTERVAL_MILLISECONDS",
	)

	jaegerEndpoint := config.GetStringRequired("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := config.GetStringRequired("SERVICE_NAME")
	serviceAddress := config.GetStringRequired("SERVICE_ADDRESS")
	metricsServiceReadTimeout := config.GetSecondsDurationRequired("METRICS_SERVICE_READ_TIMEOUT_SECONDS")
	metricsServiceWriteTimeout := config.GetSecondsDurationRequired("METRICS_SERVICE_WRITE_TIMEOUT_SECONDS")
	metricsServiceIdleTimeout := config.GetSecondsDurationRequired("METRICS_SERVICE_IDLE_TIMEOUT_SECONDS")
	metricsServiceReadHeaderTimeout := config.GetSecondsDurationRequired("METRICS_SERVICE_READ_HEADER_TIMEOUT_SECONDS")
	metricsServiceListingSuffix := config.GetStringRequired("METRICS_SERVICE_LISTING_SUFFIX")
	metricsServiceAddress := config.GetStringRequired("METRICS_SERVICE_ADDRESS")

	collectDatabaseQueryMetrics := config.GetBoolRequired("COLLECT_DATABASE_QUERY_METRICS")
	traceDatabaseQuery := config.GetBoolRequired("TRACE_DATABASE_QUERY")

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
	logger.Info("Connecting to database...")
	mongoSettings := mongodb.NewMongoSettings(
		databaseURI, databaseName, databaseCollectionName, databaseConnectionTimeout, logger,
	)
	mongoCollection, err := mongodb.Connect(ctx, databaseReconnectionCount, databaseReconnectionInterval, mongoSettings)
	if err != nil {
		logger.Error("Can't establish connection with database", zap.Error(err))
		panic(err)
	}
	logger.Info("Successfully connected to database")

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

	var repository repositories.UserRepository
	repository = mongorepository.NewUserMongoRepository(mongoCollection, databaseOperationTimeout, logger)
	if collectDatabaseQueryMetrics {
		repository = decorators.NewMetricDecorator(repository)
	}
	if traceDatabaseQuery {
		// TODO: make as config option (?)
		recordTraceErrors := true
		repository = decorators.NewTracingRepositoryDecorator(repository, tracer, recordTraceErrors)
	}

	service := services.NewUserService(repository, logger, tracer)
	controller := controllers.NewUserController(service, logger, tracer)
	server := servers.NewGRPCUserServer(controller, logger, tracer)
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
