package cviper

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type CustomViper struct {
	viper.Viper
}

func New() *CustomViper {
	return &CustomViper{
		Viper: *(viper.New()),
	}
}

func (v *CustomViper) getNoKeyPanicMessage(key string) string {
	return "Can't find required key: " + key
}

func (v *CustomViper) GetStringRequired(key string) string {
	if v.Get(key) == nil {
		panic(v.getNoKeyPanicMessage(key))
	}
	return v.GetString(key)
}

func (v *CustomViper) GetIntRequired(key string) int {
	if v.Get(key) == nil {
		panic(v.getNoKeyPanicMessage(key))
	}
	return v.GetInt(key)
}

func (v *CustomViper) GetZapcoreLogLevelRequired(key string) zapcore.Level {
	value := v.GetStringRequired(key)
	switch value {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		panic("Required one of this value for key " + key + ": 'debug', 'info', 'error'. " + value + " got.")
	}
}

func (v *CustomViper) GetSecondsDurationRequired(key string) time.Duration {
	return time.Duration(v.GetIntRequired(key)) * time.Second
}
