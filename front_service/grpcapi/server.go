package grpcapi

import (
	"context"

	"github.com/YFatMR/go_messenger/front_server/apientity"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type FrontServer struct {
	proto.UnimplementedFrontServer
	proxyController       apientity.ProxyController
	unsafeProxyController apientity.UnsafeProxyController
}

func NewFrontServer(proxyController apientity.ProxyController,
	unsafeProxyController apientity.UnsafeProxyController,
) proto.FrontServer {
	return &FrontServer{
		proxyController:       proxyController,
		unsafeProxyController: unsafeProxyController,
	}
}

// unsafeProxyController - without authorization

func (s *FrontServer) CreateUser(ctx context.Context, request *proto.CreateUserFrontRequest) (
	*proto.UserID, error,
) {
	return s.unsafeProxyController.CreateUser(ctx, request)
}

func (s *FrontServer) GenerateToken(ctx context.Context, request *proto.PublicCredential) (
	*proto.Token, error,
) {
	return s.unsafeProxyController.GenerateToken(ctx, request)
}

func (s *FrontServer) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, error,
) {
	return s.unsafeProxyController.Ping(ctx, request)
}

// proxyController - with authoriation

func (s *FrontServer) GetUserByID(ctx context.Context, request *proto.UserID) (
	*proto.UserData, error,
) {
	return s.proxyController.GetUserByID(ctx, request)
}

func (s *FrontServer) GetProgramByID(ctx context.Context, request *proto.ProgramID) (
	*proto.Program, error,
) {
	return s.proxyController.GetProgramByID(ctx, request)
}

func (s *FrontServer) CreateProgram(ctx context.Context, request *proto.ProgramSource) (
	*proto.ProgramID, error,
) {
	return s.proxyController.CreateProgram(ctx, request)
}

func (s *FrontServer) UpdateProgramSource(ctx context.Context, request *proto.UpdateProgramSourceRequest) (
	*proto.Void, error,
) {
	return s.proxyController.UpdateProgramSource(ctx, request)
}

func (s *FrontServer) RunProgram(ctx context.Context, request *proto.ProgramID) (
	*proto.Void, error,
) {
	return s.proxyController.RunProgram(ctx, request)
}

func (s *FrontServer) LintProgram(ctx context.Context, request *proto.ProgramID) (
	*proto.Void, error,
) {
	return s.proxyController.LintProgram(ctx, request)
}

func (s *FrontServer) CreateDialogWith(ctx context.Context, request *proto.UserID) (
	*proto.Dialog, error,
) {
	return s.proxyController.CreateDialogWith(ctx, request)
}

func (s *FrontServer) GetDialogs(ctx context.Context, request *proto.GetDialogsRequest) (
	*proto.GetDialogsResponse, error,
) {
	return s.proxyController.GetDialogs(ctx, request)
}
