package entity

import (
	"fmt"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type Languages string

const (
	NIL          Languages = "NIL"
	PYTHON_V_3_9 Languages = "python:3.9"
)

func LanguageFromString(language string) (Languages, error) {
	switch language {
	case "python:3.9":
		return PYTHON_V_3_9, nil
	}
	return NIL, fmt.Errorf("unsupported language")
}

type UserID struct {
	ID uint64
}

func VoidProtobuf() *proto.Void {
	return &proto.Void{}
}
