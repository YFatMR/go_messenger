package logerr

import "github.com/YFatMR/go_messenger/core/pkg/errors/logerr/logerrcore"

func Skip() logerrcore.Field {
	return logerrcore.Field{Type: logerrcore.SkipType}
}

func Int(key string, value int) logerrcore.Field {
	return Int64(key, int64(value))
}

func Int64(key string, value int64) logerrcore.Field {
	return logerrcore.Field{Key: key, Type: logerrcore.Int64Type, Integer: value}
}

func String(key string, value string) logerrcore.Field {
	return logerrcore.Field{Key: key, Type: logerrcore.StringType, String: value}
}

func Err(err error) logerrcore.Field {
	return NamedError("error", err)
}

func NamedError(key string, err error) logerrcore.Field {
	if err == nil {
		return Skip()
	}
	return logerrcore.Field{Key: key, Type: logerrcore.ErrorType, Error: err}
}
