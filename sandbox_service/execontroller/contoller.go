package execontroller

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/sandbox_service/entities/program"
	"github.com/YFatMR/go_messenger/sandbox_service/exeservice"
)

type ProgramExecutionController interface {
	Execute(ctx context.Context, request *proto.Program) (
		programResult *proto.ProgramResult, logStash ulo.LogStash, err error,
	)
	Ping(ctx context.Context, request *proto.Void) (
		pong *proto.Pong, logStash ulo.LogStash, err error,
	)
}

type programExecutionController struct {
	programExecutionService exeservice.ProgramExecutionService
}

func New(programExecutionService exeservice.ProgramExecutionService) ProgramExecutionController {
	return &programExecutionController{
		programExecutionService: programExecutionService,
	}
}

func (c *programExecutionController) Execute(ctx context.Context, request *proto.Program) (
	*proto.ProgramResult, ulo.LogStash, error,
) {
	program, err := program.FromProtobuf(request)
	if err != nil {
		return nil, ulo.FromError(err), ErrProgramExecution
	}

	programResult, _, err := c.programExecutionService.Execute(ctx, program)
	if err != nil {
		return nil, nil, err
	}

	return &proto.ProgramResult{
		Stdout: programResult.GetStdOut(),
		Stderr: programResult.GetStdErr(),
	}, nil, nil
}

func (c *programExecutionController) Ping(ctx context.Context, request *proto.Void) (
	*proto.Pong, ulo.LogStash, error,
) {
	return &proto.Pong{
		Message: "pong",
	}, nil, nil
}
