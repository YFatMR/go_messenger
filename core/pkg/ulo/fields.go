package ulo

import "github.com/YFatMR/go_messenger/core/pkg/ulo/ulocore"

func Int(key string, value int) ulocore.Field {
	return ulocore.Int(key, value)
}

func Int64(key string, value int64) ulocore.Field {
	return ulocore.Int64(key, value)
}

func String(key string, value string) ulocore.Field {
	return ulocore.String(key, value)
}

func Message(value string) ulocore.Field {
	return ulocore.Message(value)
}

func Error(err error) ulocore.Field {
	return ulocore.Error(err)
}

func NamedError(key string, err error) ulocore.Field {
	return ulocore.NamedError(key, err)
}

func Skip() ulocore.Field {
	return ulocore.Skip()
}
