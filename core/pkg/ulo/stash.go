package ulo

import (
	"github.com/YFatMR/go_messenger/core/pkg/ulo/ulocore"
	"github.com/YFatMR/go_messenger/core/pkg/ulo/ulomessage"
)

type LogStash interface {
	GetMessages() []*ulomessage.Message
}

type logstash struct {
	messages []*ulomessage.Message
}

func New(messages ...*ulomessage.Message) *logstash {
	return &logstash{
		messages: messages,
	}
}

func ErrorMsg(fields ...ulocore.Field) *logstash {
	messages := []*ulomessage.Message{ulomessage.New(ulocore.ErrorLevel, fields...)}
	return &logstash{
		messages: messages,
	}
}

func (e *logstash) GetMessages() []*ulomessage.Message {
	if e == nil {
		return []*ulomessage.Message{}
	}
	return e.messages
}
