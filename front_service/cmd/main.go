package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/ctrace"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/front_server/httpapi"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/exporters/jaeger"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := cviper.New()
	config.AutomaticEnv()

	// Init environment vars
	restFrontServiceAddress := config.GetStringRequired("REST_SERVICE_ADDRESS")
	// grpcFrontServiceAddress := config.GetStringRequired("GRPC_SERVICE_ADDRESS")
	// dialogServiceAddress := config.GetStringRequired("DIALOG_SERVICE_ADDRESS")
	jaegerEndpoint := config.GetStringRequired("JAEGER_COLLECTOR_ENDPOINT")
	serviceName := config.GetStringRequired("SERVICE_NAME")
	restServiceReadTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_READ_TIMEOUT_SECONDS")
	restServiceWriteTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_WRITE_TIMEOUT_SECONDS")
	restServiceIdleTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_IDLE_TIMEOUT_SECONDS")
	restServiceReadHeaderTimeout := config.GetSecondsDurationRequired("REST_FRONT_SERVICE_READ_HEADER_TIMEOUT_SECONDS")

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
	defer func(ctx context.Context) {
		if err := traceProvider.Shutdown(ctx); err != nil { // TODO: Shutdown -> with timeout
			panic(err)
		}
	}(ctx)
	// tracer := otel.Tracer(serviceName)

	userServiceClient, err := NewGRPCUserServiceClientFromConfig(ctx, config, logger)
	if err != nil {
		panic(err)
	}
	sandboxServiceClient, err := NewGRPCSandboxServiceClientFromConfig(ctx, config, logger)
	if err != nil {
		panic(err)
	}

	dialogServiceClient, err := NewGRPCDialogServiceClientFromConfig(ctx, config, logger)
	if err != nil {
		panic(err)
	}

	server := httpapi.NewFrontServer(userServiceClient, sandboxServiceClient, dialogServiceClient)
	router := mux.NewRouter()

	router.HandleFunc("/v1/users/{ID}", server.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/v1/token", server.GenerateToken).Methods(http.MethodPost)
	router.HandleFunc("/v1/users", server.CreateUser).Methods(http.MethodPost)

	router.HandleFunc("/v1/programs/{ID}", server.GetProgramByID).Methods(http.MethodGet)
	router.HandleFunc("/v1/programs", server.CreateProgram).Methods(http.MethodPost)
	router.HandleFunc("/v1/programs/{ID}/run", server.RunProgram).Methods(http.MethodPatch)
	router.HandleFunc("/v1/programs/{ID}/lint", server.LintProgram).Methods(http.MethodPatch)
	router.HandleFunc("/v1/programs/{ID}/source", server.LintProgram).Methods(http.MethodPatch)

	router.HandleFunc("/v1/dialogs", server.CreateDialogWith).Methods(http.MethodPost)
	router.HandleFunc("/v1/dialogs", server.GetDialogs).Methods(http.MethodGet)
	router.HandleFunc("/v1/dialogs/{ID}/messages", server.CreateDialogMessage).Methods(http.MethodPost)
	router.HandleFunc("/v1/dialogs/{ID}/messages", server.GetDialogMessages).Methods(http.MethodGet)

	service := http.Server{
		ReadTimeout:       restServiceReadTimeout,
		WriteTimeout:      restServiceWriteTimeout,
		IdleTimeout:       restServiceIdleTimeout,
		ReadHeaderTimeout: restServiceReadHeaderTimeout,
		Addr:              restFrontServiceAddress,
		Handler:           CORSMiddleware(router),
	}
	logger.Info(
		"Starting to register REST user server",
		zap.String("REST front server address", restFrontServiceAddress),
	)
	if err := service.ListenAndServe(); err != nil {
		panic(err)
	}

	// Run REST front service with "user service" API
	// go func() {
	// 	mux := runtime.NewServeMux()

	// 	service := http.Server{
	// 		ReadTimeout:       restServiceReadTimeout,
	// 		WriteTimeout:      restServiceWriteTimeout,
	// 		IdleTimeout:       restServiceIdleTimeout,
	// 		ReadHeaderTimeout: restServiceReadHeaderTimeout,
	// 		Addr:              restFrontServiceAddress,
	// 		Handler:           enableCors(mux),
	// 	}
	// 	logger.Info(
	// 		"Starting to register REST user server",
	// 		zap.String("REST front server address", restFrontServiceAddress),
	// 		zap.String("gRPC front service address", grpcFrontServiceAddress),
	// 	)

	// 	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// 	if err := proto.RegisterFrontHandlerFromEndpoint(ctx, mux, grpcFrontServiceAddress, opts); err != nil {
	// 		panic(err)
	// 	}

	// 	logger.Info("Starting serve REST front server")
	// 	//#nosec G114: Use of net/http serve function that has no support for setting timeouts
	// 	if err := service.ListenAndServe(); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// // Run gRPC front service with "user service" API
	// func() {
	// 	grpcServer := grpc.NewServer()

	// 	listener, err := net.Listen("tcp", grpcFrontServiceAddress)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	logger.Info(
	// 		"Starting to register gRPC user server",
	// 		zap.String("grpc front server address", grpcFrontServiceAddress),
	// 		zap.String("user server address", userServiceAddress),
	// 	)

	// 	logger.Info("Connecting to user service...", zap.String("address", userServiceAddress))

	// 	headers := grpcapi.GRPCHeadersFromConfig(config)
	// 	jwtManager := jwtmanager.FromConfig(config, logger)

	// 	unaryInterceptors := []grpc.UnaryClientInterceptor{
	// 		otelgrpc.UnaryClientInterceptor(),
	// 		grpcapi.UnaryAuthInterceptor(jwtManager, headers, logger),
	// 	}

	// 	grpcClientOpts := []grpc.DialOption{
	// 		grpc.WithKeepaliveParams(grpcKeepaliveParameters),
	// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 		grpc.WithUnaryInterceptor(
	// 			middleware.ChainUnaryClient(
	// 				unaryInterceptors...,
	// 			),
	// 		),
	// 		grpc.WithConnectParams(
	// 			grpc.ConnectParams{
	// 				Backoff: grpcBackoffConfig,
	// 			},
	// 		),
	// 	}
	// 	userServiceClient, err := grpcclients.NewGRPCUserClient(
	// 		ctx, userServiceAddress, microservicesConnectionTimeout, grpcClientOpts,
	// 	)
	// 	if err != nil {
	// 		logger.Error("Server stopped! Can't connect to user service", zap.String("address", userServiceAddress))
	// 		panic(err)
	// 	}

	// 	sandboxServiceClient, err := grpcclients.NewGRPCSandboxClient(
	// 		ctx, sandboxServiceAddress, microservicesConnectionTimeout, grpcClientOpts,
	// 	)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	dialogServiceClient, err := grpcclients.NewGRPCDialogClient(
	// 		ctx, dialogServiceAddress, microservicesConnectionTimeout, grpcClientOpts,
	// 	)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	proxyController := proxy.NewController(userServiceClient, sandboxServiceClient, dialogServiceClient, logger)
	// 	unsafeProxyController := proxy.NewUnsafeController(userServiceClient, logger)

	// 	server := grpcapi.NewFrontServer(proxyController, unsafeProxyController)
	// 	proto.RegisterFrontServer(grpcServer, server)

	// 	logger.Info("Starting serve gRPC front server")
	// 	if err := grpcServer.Serve(listener); err != nil {
	// 		panic(err)
	// 	}
	// }()
}
