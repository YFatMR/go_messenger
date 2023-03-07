package entity

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type ProgramSource struct {
	SourceCode string
	Language   string //  TODO: Enum
}

func ProgramSourceFromProtobuf(programSource *proto.ProgramSource) (
	*ProgramSource, error,
) {
	if programSource == nil {
		return nil, ErrWrongRequestFormat
	}
	return &ProgramSource{
		Language:   programSource.GetLanguage(),
		SourceCode: programSource.GetSourceCode(),
	}, nil
}

func ProgramSourceToProtobuf(programSource *ProgramSource) *proto.ProgramSource {
	if programSource == nil {
		return &proto.ProgramSource{}
	}
	return &proto.ProgramSource{
		Language:   programSource.Language,
		SourceCode: programSource.SourceCode,
	}
}
