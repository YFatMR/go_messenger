package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/sandbox_service/entity"
)

type SandboxRepository interface {
	GetProgramByID(ctx context.Context, programID *entity.ProgramID) (
		program *entity.Program, err error,
	)
	CreateProgram(ctx context.Context, program *entity.ProgramSource) (
		programID *entity.ProgramID, err error,
	)
	UpdateProgramSource(ctx context.Context, programID *entity.ProgramID, programSource *entity.ProgramSource) (
		err error,
	)
	UpdateCodeRunnerOutput(ctx context.Context, programID *entity.ProgramID, newOutput *entity.ProgramOutput) (
		err error,
	)
	UpdateLinterOutput(ctx context.Context, programID *entity.ProgramID, newOutput *entity.ProgramOutput) (
		err error,
	)
}
