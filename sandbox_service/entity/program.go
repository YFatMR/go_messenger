package entity

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type Program struct {
	ID               ProgramID
	Source           ProgramSource
	CodeRunnerOutput ProgramOutput
	LinterOutput     ProgramOutput
}

func ProgramToProtobuf(program *Program) *proto.Program {
	if program == nil {
		return &proto.Program{}
	}
	return &proto.Program{
		ProgramID:        ProgramIDToProtobuf(&program.ID),
		Source:           ProgramSourceToProtobuf(&program.Source),
		CodeRunnerOutput: ProgramOutputToProtobuf(&program.CodeRunnerOutput),
		LinterOutput:     ProgramOutputToProtobuf(&program.LinterOutput),
	}
}
