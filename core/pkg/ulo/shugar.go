package ulo

import (
	"github.com/YFatMR/go_messenger/core/pkg/ulo/ulocore"
	"github.com/YFatMR/go_messenger/core/pkg/ulo/ulolog"
)

func FromErrorMsg(message string, fields ...ulocore.Field) LogStash {
	newFields := make([]ulocore.Field, 0, len(fields)+1)
	newFields = append(newFields, Message(message))
	newFields = append(newFields, fields...)

	messages := []*ulolog.Log{ulolog.New(ulocore.ErrorLevel, newFields...)}
	return &logstash{
		messages: messages,
	}
}

func FromError(err error, fields ...ulocore.Field) LogStash {
	newFields := make([]ulocore.Field, 0, len(fields)+1)
	newFields = append(newFields, Error(err))
	newFields = append(newFields, fields...)

	messages := []*ulolog.Log{ulolog.New(ulocore.ErrorLevel, newFields...)}
	return &logstash{
		messages: messages,
	}
}

func FromErrorWithMsg(message string, err error, fields ...ulocore.Field) LogStash {
	return FromError(err, Message(message))
}
