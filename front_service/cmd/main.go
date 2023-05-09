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

	websocketClient := WebsocketClientFromConfig(config, logger)

	server := httpapi.NewFrontServer(
		userServiceClient, sandboxServiceClient, dialogServiceClient, websocketClient, logger,
	)
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
	router.HandleFunc("/v1/dialogs/{ID}", server.GetDialogByID).Methods(http.MethodGet)
	router.HandleFunc("/v1/dialogs", server.GetDialogs).Methods(http.MethodGet)
	router.HandleFunc("/v1/dialogs/{ID}/messages", server.CreateDialogMessage).Methods(http.MethodPost)
	router.HandleFunc("/v1/dialogs/{dialogID}/messages/{messageID}", server.GetDialogMessages).Methods(http.MethodGet)
	router.HandleFunc("/v1/dialogs/{dialogID}/messages/{messageID}", server.ReadAllMessagesBefore).Methods(http.MethodPut)

	router.HandleFunc("/v1/ws", server.WebsocketHandler)

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
}
