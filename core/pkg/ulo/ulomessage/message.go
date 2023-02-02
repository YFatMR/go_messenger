package ulomessage

import "github.com/YFatMR/go_messenger/core/pkg/ulo/ulocore"

type Message struct {
	logLevel ulocore.LogLevel
	fields   []ulocore.Field
}

func New(logLevel ulocore.LogLevel, fields ...ulocore.Field) *Message {
	return &Message{logLevel: logLevel, fields: fields}
}

func Debug(fields ...ulocore.Field) *Message {
	return &Message{logLevel: ulocore.DebugLevel, fields: fields}
}

func Info(fields ...ulocore.Field) *Message {
	return &Message{logLevel: ulocore.InfoLevel, fields: fields}
}

func Warning(fields ...ulocore.Field) *Message {
	return &Message{logLevel: ulocore.WarningLevel, fields: fields}
}

func Error(fields ...ulocore.Field) *Message {
	return &Message{logLevel: ulocore.ErrorLevel, fields: fields}
}

func (m *Message) GetLogLevel() ulocore.LogLevel {
	if m == nil {
		return ulocore.UnknownLevel
	}
	return m.logLevel
}

func (m *Message) GetFields() []ulocore.Field {
	if m == nil {
		return []ulocore.Field{}
	}
	return m.fields
}
