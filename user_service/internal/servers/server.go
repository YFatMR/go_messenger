package servers

import (
	"go.uber.org/zap"
	"golang.org/x/net/context"
	proto "protocol/pkg/proto"
)

// // 1 -> many
// // DBCLient -> cinnection
// // Controller -> create_user, ...
// // Repository   -> call create_user 2 times for example (buisness logic)
// // Server (Endpoints) ->

type userController interface {
	Create(ctx context.Context, request *proto.UserData) (*proto.UserId, error)
	GetById(ctx context.Context, request *proto.UserId) (*proto.UserData, error)
}

type GRPCUserServer struct {
	proto.UnimplementedUserServer
	controller userController
	logger     *zap.Logger
}

func NewGRPCUserServer(controller userController) *GRPCUserServer {
	return &GRPCUserServer{
		controller: controller,
	}
}

func (s *GRPCUserServer) CreateUser(ctx context.Context, request *proto.UserData) (*proto.UserId, error) {
	return s.controller.Create(ctx, request)
}

func (s *GRPCUserServer) GetUserById(ctx context.Context, request *proto.UserId) (*proto.UserData, error) {
	return s.controller.GetById(ctx, request)
}
