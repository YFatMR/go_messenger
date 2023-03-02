package protobufentity

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/sandbox_service/entity"
)

func ToProgramID(programID *proto.ProgramID) (*entity.ProgramID, error) {
	if programID == nil || programID.GetID() == "" {
		return nil, ErrWrongRequestFormat
	}
	return &entity.ProgramID{
		ID: programID.GetID(),
	}, nil
}

func FromProgramID(programID *entity.ProgramID) *proto.ProgramID {
	if programID == nil {
		return &proto.ProgramID{}
	}
	return &proto.ProgramID{
		ID: programID.ID,
	}
}
