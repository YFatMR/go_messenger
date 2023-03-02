package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

// ProxyController provide API for handlers with authorization.
type ProxyController interface {
	GetUserByID(ctx context.Context, request *proto.UserID) (
		userData *proto.UserData, err error,
	)
	GetProgramByID(ctx context.Context, request *proto.ProgramID) (
		program *proto.Program, err error,
	)
	CreateProgram(ctx context.Context, request *proto.ProgramSource) (
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
}
