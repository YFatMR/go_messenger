package program

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/YFatMR/go_messenger/sandbox_service/entities"
)

type Entity struct {
	sourceCode string
	language   string //  TODO: Enum
}

func New(language string, sourceCode string) *Entity {
	return &Entity{
		sourceCode: sourceCode,
		language:   language,
	}
}

func FromProtobuf(program *proto.Program) (*Entity, error) {
	if program == nil || program.GetLanguage() == "" || program.GetSourceCode() == "" {
		return nil, entities.ErrWrongRequestFormat
	}
	return New(program.GetLanguage(), program.GetSourceCode()), nil
}

func (e *Entity) GetSourceCode() string {
	return e.sourceCode
}

func (e *Entity) GetLanguage() string {
	return e.language
}
