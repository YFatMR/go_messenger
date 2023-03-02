package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type SandboxController interface {
	GetProgramByID(ctx context.Context, request *proto.ProgramID) (
		program *proto.Program, err error,
	)
	CreateProgram(context.Context, *proto.ProgramSource) (
		programID *proto.ProgramID, err error,
	)
	UpdateProgramSource(ctx context.Context, request *proto.UpdateProgramSourceRequest) (
		void *proto.Void, err error,
	)
	RunProgram(ctx context.Context, request *proto.ProgramID) (
		void *proto.Void, err error,
	)
	LintProgram(ctx context.Context, request *proto.ProgramID) (
		void *proto.Void, err error,
	)
	Ping(ctx context.Context, request *proto.Void) (
		pong *proto.Pong, err error,
	)
}
