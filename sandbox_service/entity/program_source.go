package entity

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type ProgramSource struct {
	SourceCode string
	Language   Languages
}

func ProgramSourceFromProtobuf(programSource *proto.ProgramSource) (
	*ProgramSource, error,
) {
	if programSource == nil {
		return nil, ErrWrongRequestFormat
	}
	language, err := LanguageFromString(programSource.GetLanguage())
	if err != nil {
		return nil, err
	}

	return &ProgramSource{
		Language:   language,
		SourceCode: programSource.GetSourceCode(),
	}, nil
}

func ProgramSourceToProtobuf(programSource *ProgramSource) *proto.ProgramSource {
	if programSource == nil {
		return &proto.ProgramSource{}
	}
	return &proto.ProgramSource{
		Language:   string(programSource.Language),
		SourceCode: programSource.SourceCode,
	}
}
