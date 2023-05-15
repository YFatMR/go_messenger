package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/ctrace"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/sandbox_service/decorator"
	"github.com/YFatMR/go_messenger/sandbox_service/grpcapi"
	"github.com/YFatMR/go_messenger/sandbox_service/sandbox"
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

	serviceAddress := config.GetStringRequired("SERVICE_ADDRESS")
	serviceName := config.GetStringRequired("SERVICE_NAME")
	jaegerEndpoint := config.GetStringRequired("JAEGER_COLLECTOR_ENDPOINT")

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

	sandboxRepository, err := sandbox.RepositoryFromConfig(ctx, config, logger)
	if err != nil {
		panic(err)
	}
	sandboxRepository = decorator.NewOpentelemetryTracingSandboxRepositoryDecorator(
		sandboxRepository, tracer, false,
	)
	sandboxRepository = decorator.NewLoggingSandboxRepositoryDecorator(sandboxRepository, logger)

	kafkaClient := sandbox.KafkaClientFromConfig(config, logger)
	defer kafkaClient.Stop()
	sandboxService := sandbox.NewService(sandboxRepository, kafkaClient, logger)
	sandboxService = decorator.NewOpentelemetryTracingSandboxServiceDecorator(sandboxService, tracer, false)
	sandboxService = decorator.NewLoggingSandboxServiceDecorator(sandboxService, logger)

	grpcHeaders := grpcapi.HeadersFromConfig(config)
	contextManager := grpcapi.NewContextManager(grpcHeaders)
	controller := sandbox.NewController(sandboxService, contextManager, logger)
	controller = decorator.NewLoggingSandboxControllerDecorator(controller, logger)

	sandboxServer := grpcapi.NewSandboxServer(controller)

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		panic(err)
	}
	proto.RegisterSandboxServer(grpcServer, &sandboxServer)
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
