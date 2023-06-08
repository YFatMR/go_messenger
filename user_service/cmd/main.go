package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/ctrace"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/metrics/prometheus"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/grpcapi"
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
	metricServiceSettings := cviper.NewMetricMetricServiceSettingsFromConfig(config)

	jaegerEndpoint := config.GetStringRequired("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := config.GetStringRequired("SERVICE_NAME")
	serviceAddress := config.GetStringRequired("SERVICE_ADDRESS")

	// Init logger
	logger, err := czap.FromConfig(config)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic!", zap.Any("msg", r))
		}
		panic("Panic")
	}()

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
	repository, err := UserRepositoryFromConfig(ctx, config, logger, tracer)
	if err != nil {
		panic(err)
	}
	service := UserServiceFromConfig(config, logger, repository, tracer)
	controller := UserControllerFromConfig(config, service, logger)

	server := grpcapi.NewServer(controller)
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
