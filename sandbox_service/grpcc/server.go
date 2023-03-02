package grpcc

import (
	"context"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/sandbox_service/apientity"
)

type SandboxServer struct {
	proto.UnimplementedSandboxServer
	sandboxController apientity.SandboxController
}

func NewSandboxServer(sandboxController apientity.SandboxController) SandboxServer {
	return SandboxServer{
		sandboxController: sandboxController,
	}
}

func (s *SandboxServer) GetProgramByID(ctx context.Context, request *proto.ProgramID) (
	*proto.Program, error,
) {
	return s.sandboxController.GetProgramByID(ctx, request)
}

func (s *SandboxServer) CreateProgram(ctx context.Context, request *proto.ProgramSource) (
	*proto.ProgramID, error,
) {
	return s.sandboxController.CreateProgram(ctx, request)
}

func (s *SandboxServer) UpdateProgramSource(ctx context.Context, request *proto.UpdateProgramSourceRequest) (
	*proto.Void, error,
) {
	return s.sandboxController.UpdateProgramSource(ctx, request)
}

func (s *SandboxServer) RunProgram(ctx context.Context, request *proto.ProgramID) (
	*proto.Void, error,
) {
	return s.sandboxController.RunProgram(ctx, request)
}

func (s *SandboxServer) LintProgram(ctx context.Context, request *proto.ProgramID) (
	*proto.Void, error,
) {
	return s.sandboxController.LintProgram(ctx, request)
}

func (s *SandboxServer) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, error,
) {
	return s.sandboxController.Ping(ctx, request)
}
