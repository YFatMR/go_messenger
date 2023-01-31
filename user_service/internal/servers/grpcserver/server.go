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

func (s *GRPCServer) CreateUser(ctx context.Context, request *proto.CreateUserDataRequest) (*proto.UserID, error) {
	userID, lerr := s.controller.Create(ctx, request)
	return userID, lerr.GetAPIError()
}

func (s *GRPCServer) GetUserByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error) {
	userData, lerr := s.controller.GetByID(ctx, request)
	return userData, lerr.GetAPIError()
}

func (s *GRPCServer) DeleteUserByID(ctx context.Context, request *proto.UserID) (*proto.Void, error) {
	void, lerr := s.controller.DeleteByID(ctx, request)
	return void, lerr.GetAPIError()
}

func (s *GRPCServer) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, error) {
	pong, lerr := s.controller.Ping(ctx, request)
	return pong, lerr.GetAPIError()
}
