package grpcserver

import (
	"context"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/sandbox_service/execontroller"
)

type ProgramExecution struct {
	proto.UnimplementedSandboxServer
	programExecutionController execontroller.ProgramExecutionController
}

func New(programExecutionController execontroller.ProgramExecutionController) ProgramExecution {
	return ProgramExecution{
		programExecutionController: programExecutionController,
	}
}

func (s *ProgramExecution) Execute(ctx context.Context, request *proto.Program) (
	*proto.ProgramResult, error,
) {
	programResult, _, err := s.programExecutionController.Execute(ctx, request)
	return programResult, err
}

func (s *ProgramExecution) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, error,
) {
	pong, _, err := s.programExecutionController.Ping(ctx, request)
	return pong, err
}
