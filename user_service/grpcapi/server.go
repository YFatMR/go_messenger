package grpcapi

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/user_service/apientity"
	"golang.org/x/net/context"
)

type Server struct {
	proto.UnimplementedUserServer
	controller apientity.UserController
}

func NewServer(controller apientity.UserController) Server {
	return Server{
		controller: controller,
	}
}

func (s *Server) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (
	*proto.UserID, error,
) {
	return s.controller.Create(ctx, request)
}

func (s *Server) GetUserByID(ctx context.Context, request *proto.UserID) (*proto.UserData, error) {
	return s.controller.GetByID(ctx, request)
}

func (s *Server) DeleteUserByID(ctx context.Context, request *proto.UserID) (*proto.Void, error) {
	return s.controller.DeleteByID(ctx, request)
}

func (s *Server) GenerateToken(ctx context.Context, request *proto.Credential) (*proto.Token, error) {
	return s.controller.GenerateToken(ctx, request)
}

func (s *Server) Ping(ctx context.Context, request *proto.Void) (*proto.Pong, error) {
	return s.controller.Ping(ctx, request)
}
