package main

import (
	"context"
	"net/http"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/grpcclients"
	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
	"github.com/YFatMR/go_messenger/front_server/grpcapi"
	"github.com/YFatMR/go_messenger/front_server/websocketapi"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/gorilla/websocket"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func CORSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func KeepaliveParamsFromConfig(config *cviper.CustomViper) keepalive.ClientParameters {
	return keepalive.ClientParameters{
		Time:                config.GetMillisecondsDurationRequired("GRPC_CONNECTION_KEEPALIVE_TIME_MILLISECONDS"),
		Timeout:             config.GetMillisecondsDurationRequired("GRPC_CONNECTION_KEEPALIVE_TIMEOUT_MILLISECONDS"),
		PermitWithoutStream: config.GetBoolRequired("GRPC_CONNECTION_KEEPALIVE_PERMIT_WITHOUT_STREAM"),
	}
}

func BackoffSettingsFromConfig(config *cviper.CustomViper) backoff.Config {
	return backoff.Config{
		BaseDelay:  config.GetMillisecondsDurationRequired("GRPC_CONNECTION_BACKOFF_DELAY_MILLISECONDS"),
		Multiplier: config.GetFloat64Required("GRPC_CONNECTION_BACKOFF_MULTIPLIER"),
		Jitter:     config.GetFloat64Required("GRPC_CONNECTION_BACKOFF_JITTER"),
		MaxDelay:   config.GetMillisecondsDurationRequired("GRPC_CONNECTION_BACKOFF_MAX_DELAY_MILLISECONDS"),
	}
}

func DefaultGRPCClientOptsFromConfig(config *cviper.CustomViper, logger *czap.Logger) []grpc.DialOption {
	headers := grpcapi.GRPCHeadersFromConfig(config)
	jwtManager := jwtmanager.FromConfig(config, logger)

	unaryInterceptors := []grpc.UnaryClientInterceptor{
		otelgrpc.UnaryClientInterceptor(),
		grpcapi.UnaryAuthInterceptor(jwtManager, headers, logger),
	}

	grpcKeepaliveParameters := KeepaliveParamsFromConfig(config)
	grpcBackoffConfig := BackoffSettingsFromConfig(config)
	return []grpc.DialOption{
		grpc.WithKeepaliveParams(grpcKeepaliveParameters),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			middleware.ChainUnaryClient(
				unaryInterceptors...,
			),
		),
		grpc.WithConnectParams(
			grpc.ConnectParams{
				Backoff: grpcBackoffConfig,
			},
		),
	}
}

func NewGRPCUserServiceClientFromConfig(ctx context.Context, config *cviper.CustomViper, logger *czap.Logger) (
	proto.UserClient, error,
) {
	grpcOpts := DefaultGRPCClientOptsFromConfig(config, logger)
	userServiceAddress := config.GetStringRequired("USER_SERVICE_ADDRESS")
	microservicesConnectionTimeout := config.GetMillisecondsDurationRequired(
		"MICROSERVICES_GRPC_CONNECTION_TIMEOUT_MILLISECONDS",
	)
	return grpcclients.NewGRPCUserClient(
		ctx, userServiceAddress, microservicesConnectionTimeout, grpcOpts,
	)
}

func NewGRPCSandboxServiceClientFromConfig(ctx context.Context, config *cviper.CustomViper, logger *czap.Logger) (
	proto.SandboxClient, error,
) {
	grpcOpts := DefaultGRPCClientOptsFromConfig(config, logger)
	sandboxServiceAddress := config.GetStringRequired("SANDBOX_SERVICE_ADDRESS")
	microservicesConnectionTimeout := config.GetMillisecondsDurationRequired(
		"MICROSERVICES_GRPC_CONNECTION_TIMEOUT_MILLISECONDS",
	)
	return grpcclients.NewGRPCSandboxClient(
		ctx, sandboxServiceAddress, microservicesConnectionTimeout, grpcOpts,
	)
}

func NewGRPCDialogServiceClientFromConfig(ctx context.Context, config *cviper.CustomViper, logger *czap.Logger) (
	proto.DialogServiceClient, error,
) {
	grpcOpts := DefaultGRPCClientOptsFromConfig(config, logger)
	dialogServiceAddress := config.GetStringRequired("DIALOG_SERVICE_ADDRESS")
	microservicesConnectionTimeout := config.GetMillisecondsDurationRequired(
		"MICROSERVICES_GRPC_CONNECTION_TIMEOUT_MILLISECONDS",
	)
	return grpcclients.NewGRPCDialogClient(
		ctx, dialogServiceAddress, microservicesConnectionTimeout, grpcOpts,
	)
}

func NewGRPCBotsServiceClientFromConfig(ctx context.Context, config *cviper.CustomViper, logger *czap.Logger) (
	proto.BotsServiceClient, error,
) {
	grpcOpts := DefaultGRPCClientOptsFromConfig(config, logger)
	dialogServiceAddress := config.GetStringRequired("BOTS_SERVICE_ADDRESS")
	microservicesConnectionTimeout := config.GetMillisecondsDurationRequired(
		"MICROSERVICES_GRPC_CONNECTION_TIMEOUT_MILLISECONDS",
	)
	return grpcclients.NewGRPCBotsClient(
		ctx, dialogServiceAddress, microservicesConnectionTimeout, grpcOpts,
	)
}

func WebsocketClientFromConfig(config *cviper.CustomViper, logger *czap.Logger) *websocketapi.Client {
	return websocketapi.NewClient(
		jwtmanager.FromConfig(config, logger),
		websocket.DefaultDialer,
		&websocketapi.ClientSettings{
			Addr: config.GetStringRequired("COMET_SERVICE_WS_ADDRESS"),
		},
		&websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		logger,
	)
}
