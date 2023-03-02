package protobufentity

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/sandbox_service/entity"
)

func ToProgramSource(programSource *proto.ProgramSource) (
	*entity.ProgramSource, error,
) {
	if programSource == nil {
		return nil, ErrWrongRequestFormat
	}
	return &entity.ProgramSource{
		Language:   programSource.GetLanguage(),
		SourceCode: programSource.GetSourceCode(),
	}, nil
}

func FromProgramSource(programSource *entity.ProgramSource) *proto.ProgramSource {
	if programSource == nil {
		return &proto.ProgramSource{}
	}
	return &proto.ProgramSource{
		Language:   programSource.Language,
		SourceCode: programSource.SourceCode,
	}
}
