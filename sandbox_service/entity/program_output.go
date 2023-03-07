package entity

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type ProgramOutput struct {
	Stdout string
	Stderr string
}

func ProgramOutputToProtobuf(programOutput *ProgramOutput) *proto.ProgramOutput {
	if programOutput == nil {
		return &proto.ProgramOutput{}
	}
	return &proto.ProgramOutput{
		Stdout: programOutput.Stdout,
		Stderr: programOutput.Stderr,
	}
}
