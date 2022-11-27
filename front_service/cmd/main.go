package main

import (
	"context"
	. "core/pkg/utils"
	"fmt"
	. "front_server/internal"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

func runRest(ctx context.Context, mux *runtime.ServeMux, restFrontServiceAddress string, grpcFrontServiceAddress string) {
	time.Sleep(time.Second * 2)
	RegisterRestUserServer(ctx, mux, grpcFrontServiceAddress)
	fmt.Println("serve rest", restFrontServiceAddress)
	if err := http.ListenAndServe(restFrontServiceAddress, mux); err != nil {
		panic(err)
	}
}

func runGrpc(grpcServer *grpc.Server, grpcFrontServiceAddress string, userServerAddress string) {
	listener, err := net.Listen("tcp", grpcFrontServiceAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	RegisterGrpcUserServer(grpcServer, userServerAddress)
	fmt.Println("serve grpc")
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func main() {
	// Init vars
	frontRestUserServerAddress := GetFullServiceAddress("FRONT_REST") // :8000
	frontGrpcUserServerAddress := GetFullServiceAddress("FRONT_GRPC") // :8000
	userServerAddress := GetFullServiceAddress("USER")

	mux := runtime.NewServeMux()
	grpcServer := grpc.NewServer()
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second) // TODO: set timeout throwout config
	//defer cancel()

	go runRest(ctx, mux, frontRestUserServerAddress, frontGrpcUserServerAddress)
	runGrpc(grpcServer, frontGrpcUserServerAddress, userServerAddress)
}
