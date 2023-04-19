package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/dialog_service/grpcapi"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
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
	repository, err := DialogRepositoryFromConfig(ctx, config, logger)
	if err != nil {
		panic(err)
	}

	model, err := DialogModelFromConfig(ctx, repository, config, logger)
	if err != nil {
		panic(err)
	}

	controller, err := DialogControllerFromConfig(model, config, logger)
	if err != nil {
		panic(err)
	}

	server := grpcapi.NewServer(controller)
	s := grpc.NewServer()

	// Register protobuf server
	logger.Info("Register protobuf server")
	proto.RegisterDialogServiceServer(s, &server)

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
