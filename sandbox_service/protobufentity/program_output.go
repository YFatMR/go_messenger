package protobufentity

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/sandbox_service/entity"
)

func FromProgramOutput(programOutput *entity.ProgramOutput) *proto.ProgramOutput {
	if programOutput == nil {
		return &proto.ProgramOutput{}
	}
	return &proto.ProgramOutput{
		Stdout: programOutput.Stdout,
		Stderr: programOutput.Stderr,
	}
}
