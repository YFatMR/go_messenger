package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	userserver "github.com/YFatMR/go_messenger/front_server/internal/user_server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func runRest(ctx context.Context, mux *runtime.ServeMux, restFrontServiceAddress string,
	grpcFrontServiceAddress string, restServiceReadTimeout time.Duration,
	restServiceWriteTimeout time.Duration, logger *loggers.OtelZapLoggerWithTraceID,
) {
	service := http.Server{
		ReadTimeout:  restServiceReadTimeout,
		WriteTimeout: restServiceWriteTimeout,
		Addr:         restFrontServiceAddress,
		Handler:      mux,
	}
	logger.Info(
		"Starting to register REST user server",
		zap.String("REST front server address", restFrontServiceAddress),
		zap.String("gRPC front service address", grpcFrontServiceAddress),
	)
	userserver.RegisterRestUserServer(ctx, mux, grpcFrontServiceAddress)
	logger.Info("Starting serve REST front server")
	//#nosec G114: Use of net/http serve function that has no support for setting timeouts
	if err := service.ListenAndServe(); err != nil {
		panic(err)
	}
}

func runGrpc(grpcServer *grpc.Server, grpcFrontServiceAddress string,
	userServerAddress string, logger *loggers.OtelZapLoggerWithTraceID, tracer trace.Tracer,
) {
	listener, err := net.Listen("tcp", grpcFrontServiceAddress)
	if err != nil {
		panic(err)
	}
	logger.Info(
		"Starting to register gRPC user server",
		zap.String("grpc front server address", grpcFrontServiceAddress),
		zap.String("user server address", userServerAddress),
	)
	userserver.RegisterGrpcUserServer(grpcServer, userServerAddress, logger, tracer)
	logger.Info("Starting serve gRPC front server")
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
