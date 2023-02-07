package zapfields

import (
	"github.com/YFatMR/go_messenger/core/pkg/ulo/ulocore"
	"github.com/YFatMR/go_messenger/core/pkg/ulo/ulolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ToZapLogLevel(logLevel ulocore.LogLevel) zapcore.Level {
	switch logLevel {
	case ulocore.DebugLevel:
		return zap.DebugLevel
	case ulocore.InfoLevel:
		return zap.InfoLevel
	case ulocore.WarningLevel:
		return zap.WarnLevel
	case ulocore.ErrorLevel:
		fallthrough
	default:
		return zap.ErrorLevel
	}
}

func ToZapType(fieldtype ulocore.FieldType) zapcore.FieldType {
	switch fieldtype {
	case ulocore.StringType:
		return zapcore.StringType
	case ulocore.ErrorType:
		return zapcore.ErrorType
	case ulocore.IntType:
		return zapcore.Int32Type
	case ulocore.Int64Type:
		return zapcore.Int64Type
	case ulocore.BoolType:
		return zapcore.BoolType
	case ulocore.SkipType:
		fallthrough
	case ulocore.UnknownType:
		fallthrough
	default:
		return zapcore.SkipType
	}
}

func ToZapFiled(field ulocore.Field) zapcore.Field {
	return zapcore.Field{
		Key:       field.Key,
		Type:      ToZapType(field.Type),
		Integer:   field.Integer,
		String:    field.String,
		Interface: field.Interface,
	}
}

func FromMessage(message *ulolog.Log) []zapcore.Field {
	result := make([]zapcore.Field, 0, len(message.GetFields()))
	for _, field := range message.GetFields() {
		result = append(result, ToZapFiled(field))
	}
	return result
}
