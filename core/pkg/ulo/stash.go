package ulo

import (
	"github.com/YFatMR/go_messenger/core/pkg/ulo/ulocore"
	"github.com/YFatMR/go_messenger/core/pkg/ulo/ulolog"
)

type LogStash interface {
	GetMessages() []*ulolog.Log
	Debug(fields ...ulocore.Field)
	Info(fields ...ulocore.Field)
	Warning(fields ...ulocore.Field)
	Error(fields ...ulocore.Field)
}

type logstash struct {
	messages []*ulolog.Log
}

func New(messages ...*ulolog.Log) LogStash {
	return &logstash{
		messages: messages,
	}
}

func (l *logstash) GetMessages() []*ulolog.Log {
	if l == nil {
		return []*ulolog.Log{}
	}
	return l.messages
}

func (l *logstash) Debug(fields ...ulocore.Field) {
	if l == nil || fields == nil {
		return
	}
	l.messages = append(l.messages, ulolog.Debug(fields...))
}

func (l *logstash) Info(fields ...ulocore.Field) {
	if l == nil || fields == nil {
		return
	}
	l.messages = append(l.messages, ulolog.Info(fields...))
}

func (l *logstash) Warning(fields ...ulocore.Field) {
	if l == nil || fields == nil {
		return
	}
	l.messages = append(l.messages, ulolog.Warning(fields...))
}

func (l *logstash) Error(fields ...ulocore.Field) {
	if l == nil || fields == nil {
		return
	}
	l.messages = append(l.messages, ulolog.Error(fields...))
}
