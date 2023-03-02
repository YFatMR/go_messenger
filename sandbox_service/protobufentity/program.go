package protobufentity

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/sandbox_service/entity"
)

func FromProgram(program *entity.Program) *proto.Program {
	if program == nil {
		return &proto.Program{}
	}
	return &proto.Program{
		ProgramID:        FromProgramID(&program.ID),
		Source:           FromProgramSource(&program.Source),
		CodeRunnerOutput: FromProgramOutput(&program.CodeRunnerOutput),
		LinterOutput:     FromProgramOutput(&program.LinterOutput),
	}
}
