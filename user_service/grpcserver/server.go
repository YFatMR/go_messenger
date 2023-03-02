package grpcserver

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/controllers"
	"golang.org/x/net/context"
)

type GRPCServer struct {
	proto.UnimplementedUserServer
	controller controllers.UserController
}

func New(controller controllers.UserController) GRPCServer {
	return GRPCServer{
		controller: controller,
	}
}

func (s *GRPCServer) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (
	*proto.UserID, error,
) {
	return s.controller.Create(ctx, request)
}

func (s *GRPCServer) GetUserByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error) {
	return s.controller.GetByID(ctx, request)
}

func (s *GRPCServer) DeleteUserByID(ctx context.Context, request *proto.UserID) (*proto.Void, error) {
	return s.controller.DeleteByID(ctx, request)
}

func (s *GRPCServer) GenerateToken(ctx context.Context, request *proto.Credential) (*proto.Token, error) {
	return s.controller.GenerateToken(ctx, request)
}

func (s *GRPCServer) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, error) {
	return s.controller.Ping(ctx, request)
}
