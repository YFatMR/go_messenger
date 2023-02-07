package ulolog

import "github.com/YFatMR/go_messenger/core/pkg/ulo/ulocore"

type Log struct {
	logLevel ulocore.LogLevel
	fields   []ulocore.Field
}

func New(logLevel ulocore.LogLevel, fields ...ulocore.Field) *Log {
	return &Log{logLevel: logLevel, fields: fields}
}

func Debug(fields ...ulocore.Field) *Log {
	return &Log{logLevel: ulocore.DebugLevel, fields: fields}
}

func Info(fields ...ulocore.Field) *Log {
	return &Log{logLevel: ulocore.InfoLevel, fields: fields}
}

func Warning(fields ...ulocore.Field) *Log {
	return &Log{logLevel: ulocore.WarningLevel, fields: fields}
}

func Error(fields ...ulocore.Field) *Log {
	return &Log{logLevel: ulocore.ErrorLevel, fields: fields}
}

func (m *Log) GetLogLevel() ulocore.LogLevel {
	if m == nil {
		return ulocore.DebugLevel
	}
	return m.logLevel
}

func (m *Log) GetFields() []ulocore.Field {
	if m == nil {
		return []ulocore.Field{}
	}
	return m.fields
}
