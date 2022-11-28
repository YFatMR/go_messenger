package main

import (
	"context"
	. "core/pkg/loggers"
	. "core/pkg/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	logLevel := RequiredZapcoreLogLevelEnv("FRONT_SERVICE_LOG_LEVEL")
	logPath := RequiredStringEnv("FRONT_SERVICE_LOG_PATH")
	logger := NewBaseFileLogger(logLevel, logPath)
	defer logger.Sync()

	// Init vars
	frontRestUserServerAddress := GetFullServiceAddress("FRONT_REST")
	frontGrpcUserServerAddress := GetFullServiceAddress("FRONT_GRPC")
	userServerAddress := GetFullServiceAddress("USER")

	mux := runtime.NewServeMux()
	grpcServer := grpc.NewServer()
	ctx := context.Background()

	go runRest(ctx, mux, frontRestUserServerAddress, frontGrpcUserServerAddress)
	runGrpc(grpcServer, frontGrpcUserServerAddress, userServerAddress, logger)
}
