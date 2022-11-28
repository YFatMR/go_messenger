package internal

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "protocol/pkg/proto"
)

type frontUserServer struct {
	proto.UnimplementedFrontUserServer
	userServerAddress string
	logger            *zap.Logger
}

func newFrontUserServer(userServerAddress string, logger *zap.Logger) *frontUserServer {
	return &frontUserServer{
		userServerAddress: userServerAddress,
		logger:            logger,
	}
}

func (s *frontUserServer) CreateUser(ctx context.Context, request *proto.UserData) (*proto.UserId, error) {
	s.logger.Debug("called CreateUser endpoint")
	conn, err := grpc.Dial(s.userServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := proto.NewUserClient(conn)
	return client.CreateUser(ctx, request)
}

func (s *frontUserServer) GetUserById(ctx context.Context, request *proto.UserId) (*proto.UserData, error) {
	s.logger.Debug("called GetUserById endpoint")
	conn, err := grpc.Dial(s.userServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := proto.NewUserClient(conn)
	return client.GetUserById(ctx, request)
}

// registration

func RegisterRestUserServer(ctx context.Context, mux *runtime.ServeMux, grpcFrontServerAddress string) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := proto.RegisterFrontUserHandlerFromEndpoint(ctx, mux, grpcFrontServerAddress, opts)
	if err != nil {
		panic(err)
	}
}

func RegisterGrpcUserServer(grpcServer grpc.ServiceRegistrar, userServerAddress string, logger *zap.Logger) {
	proto.RegisterFrontUserServer(grpcServer, newFrontUserServer(userServerAddress, logger))
}
