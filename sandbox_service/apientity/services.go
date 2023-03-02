package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/sandbox_service/entity"
)

type SandboxService interface {
	GetProgramByID(ctx context.Context, programID *entity.ProgramID) (
		program *entity.Program, err error,
	)
	CreateProgram(ctx context.Context, programSource *entity.ProgramSource) (
		programID *entity.ProgramID, err error,
	)
	UpdateProgramSource(ctx context.Context, programID *entity.ProgramID, programSource *entity.ProgramSource) (
		err error,
	)
	RunProgram(ctx context.Context, programID *entity.ProgramID, userID *entity.UserID) (
		err error,
	)
	// TODO: implement
	// LintProgram(ctx context.Context, programID *entity.ProgramID) (
	// 	err error,
	// )
}
