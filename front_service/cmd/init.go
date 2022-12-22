package main

import (
	"context"
	. "github.com/YFatMR/go_messenger/core/pkg/loggers"
	. "github.com/YFatMR/go_messenger/front_server/internal/user_server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

func runRest(ctx context.Context, mux *runtime.ServeMux, restFrontServiceAddress string, grpcFrontServiceAddress string, logger *OtelZapLoggerWithTraceID) {
	logger.Info(
		"Starting to register REST user server",
		zap.String("REST front server address", restFrontServiceAddress),
		zap.String("gRPC front service address", grpcFrontServiceAddress),
	)
	RegisterRestUserServer(ctx, mux, grpcFrontServiceAddress)
	logger.Info("Starting serve REST front server")
	if err := http.ListenAndServe(restFrontServiceAddress, mux); err != nil {
		panic(err)
	}
}

func runGrpc(grpcServer *grpc.Server, grpcFrontServiceAddress string, userServerAddress string, logger *OtelZapLoggerWithTraceID, tracer trace.Tracer) {
	listener, err := net.Listen("tcp", grpcFrontServiceAddress)
	if err != nil {
		panic(err)
	}
	logger.Info(
		"Starting to register gRPC user server",
		zap.String("grpc front server address", grpcFrontServiceAddress),
		zap.String("user server address", userServerAddress),
	)
	RegisterGrpcUserServer(grpcServer, userServerAddress, logger, tracer)
	logger.Info("Starting serve gRPC front server")
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
