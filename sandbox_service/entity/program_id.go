package entity

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type ProgramID struct {
	ID string
}

func ProgramIDFromProtobuf(programID *proto.ProgramID) (*ProgramID, error) {
	if programID == nil || programID.GetID() == "" {
		return nil, ErrWrongRequestFormat
	}
	return &ProgramID{
		ID: programID.GetID(),
	}, nil
}

func ProgramIDToProtobuf(programID *ProgramID) *proto.ProgramID {
	if programID == nil {
		return &proto.ProgramID{}
	}
	return &proto.ProgramID{
		ID: programID.ID,
	}
}
