package grpcserver

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/internal/controllers"
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

func (s *GRPCServer) CreateUser(ctx context.Context, request *proto.CreateUserDataRequest) (
	*proto.UserID, error,
) {
	userID, _, err := s.controller.Create(ctx, request)
	return userID, err
}

func (s *GRPCServer) GetUserByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error) {
	userData, _, err := s.controller.GetByID(ctx, request)
	return userData, err
}

func (s *GRPCServer) DeleteUserByID(ctx context.Context, request *proto.UserID) (*proto.Void, error) {
	void, _, err := s.controller.DeleteByID(ctx, request)
	return void, err
}

func (s *GRPCServer) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, error) {
	pong, _, err := s.controller.Ping(ctx, request)
	return pong, err
}
