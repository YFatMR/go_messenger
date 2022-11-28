package main

import (
	"context"
	"fmt"
	. "front_server/internal"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

func runRest(ctx context.Context, mux *runtime.ServeMux, restFrontServiceAddress string, grpcFrontServiceAddress string) {
	RegisterRestUserServer(ctx, mux, grpcFrontServiceAddress)
	fmt.Println("serve rest", restFrontServiceAddress)
	if err := http.ListenAndServe(restFrontServiceAddress, mux); err != nil {
		panic(err)
	}
}

func runGrpc(grpcServer *grpc.Server, grpcFrontServiceAddress string, userServerAddress string, logger *zap.Logger) {
	listener, err := net.Listen("tcp", grpcFrontServiceAddress)
	if err != nil {
		panic(err)
	}
	RegisterGrpcUserServer(grpcServer, userServerAddress, logger)
	fmt.Println("serve grpc")
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
