package cviper

import (
	"fmt"
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
	fmt.Println("Can't find required key:", key)
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

func (v *CustomViper) GetInt64Required(key string) int64 {
	if v.Get(key) == nil {
		panic(v.getNoKeyPanicMessage(key))
	}
	return v.GetInt64(key)
}

func (v *CustomViper) GetFloat64Required(key string) float64 {
	if v.Get(key) == nil {
		panic(v.getNoKeyPanicMessage(key))
	}
	return v.GetFloat64(key)
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

func (v *CustomViper) GetMillisecondsDurationRequired(key string) time.Duration {
	return time.Duration(v.GetIntRequired(key)) * time.Millisecond
}

func (v *CustomViper) GetBoolRequired(key string) bool {
	if v.Get(key) == nil {
		panic(v.getNoKeyPanicMessage(key))
	}
	return v.GetBool(key)
}
