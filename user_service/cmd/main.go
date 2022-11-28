package main

import (
	"context"
	. "core/pkg/loggers"
	. "core/pkg/utils"
	"fmt"
	"google.golang.org/grpc"
	"net"
	proto "protocol/pkg/proto"
	"time"
	"user_server/internal/controllers"
	"user_server/internal/repositories"
	"user_server/internal/servers"
	"user_server/internal/services"
)

// GOMAXPROC

func main() {
	time.Sleep(5 * time.Second)

	s := grpc.NewServer()

	logLevel := RequiredZapcoreLogLevelEnv("USER_SERVICE_LOG_LEVEL")
	logPath := RequiredStringEnv("USER_SERVICE_LOG_PATH")
	logger := NewBaseFileLogger(logLevel, logPath)
	defer logger.Sync()

	mongoUri := RequiredStringEnv("USER_SERVICE_MONGODB_URI")
	databaseName := RequiredStringEnv("USER_SERVICE_MONGODB_DATABASE_NAME")
	collectionName := RequiredStringEnv("USER_SERVICE_MONGODB_DATABASE_COLLECTION_NAME")
	connectionTimeout := RequiredIntEnv("USER_SERVICE_MONGODB_CONNECTION_TIMEOUT_SEC")
	mongoSetting := NewMongoSettings(mongoUri, databaseName, collectionName, time.Duration(connectionTimeout)*time.Second)

	mongoCtx := context.Background()
	fmt.Println("Connection to mongo...")
	mongoCollection, cancelConnection := NewMongoCollection(mongoCtx, mongoSetting)
	fmt.Println("Connected to mongo...")
	defer cancelConnection()

	userRepository := repositories.NewUserMongoRepository(mongoCollection, logger)
	userService := services.NewUserService(userRepository, logger)
	userController := controllers.NewUserController(userService, logger)
	gRPCUserServer := servers.NewGRPCUserServer(userController, logger)
	proto.RegisterUserServer(s, gRPCUserServer)

	userServerAddress := GetFullServiceAddress("USER")
	listener, err := net.Listen("tcp", userServerAddress)
	if err != nil {
		panic(err)
	}
	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
