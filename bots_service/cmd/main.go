package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/bots_service/internal/grpcapi"
	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := cviper.New()
	config.AutomaticEnv()

	// Init environment vars
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

	// Init Server
	logger.Info("Init server")
	controller := BotsControllerFromConfig(config, logger)

	server := grpcapi.NewServer(controller)
	s := grpc.NewServer()

	// Register protobuf server
	logger.Info("Register protobuf server")
	proto.RegisterBotsServiceServer(s, server)

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
