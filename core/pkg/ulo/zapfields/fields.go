package zapfields

import (
	"github.com/YFatMR/go_messenger/core/pkg/ulo/ulocore"
	"github.com/YFatMR/go_messenger/core/pkg/ulo/ulomessage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ToZapFiled(field ulocore.Field) zapcore.Field {
	switch field.Type {
	case ulocore.StringType:
		return zap.String(field.Key, field.String)
	case ulocore.ErrorType:
		return zap.NamedError(field.Key, field.Interface.(error))
	case ulocore.IntType:
		return zap.Int(field.Key, int(field.Integer))
	case ulocore.Int64Type:
		return zap.Int64(field.Key, field.Integer)
	case ulocore.BoolType:
		return zap.Bool(field.Key, bool(field.Integer == 1))
	case ulocore.SkipType:
	case ulocore.UnknownType:
	default:
		return zap.Skip()
	}
	return zap.Skip()
}

func FromMessage(message *ulomessage.Message) []zapcore.Field {
	result := make([]zapcore.Field, 0, len(message.GetFields()))
	for _, field := range message.GetFields() {
		result = append(result, ToZapFiled(field))
	}
	return result
}
