package loggers

import (
	. "core/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func RequiredZapcoreLogLevelEnv(env string) zapcore.Level {
	level := RequiredStringEnv(env)
	if level == "debug" {
		return zapcore.DebugLevel
	} else if level == "info" {
		return zapcore.InfoLevel
	} else if level == "error" {
		return zapcore.ErrorLevel
	}
	panic("Variable " + env + " has unexpected values")
}

func NewBaseFileLogger(logLevel zapcore.LevelEnabler, logFilePath string) *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)

	file, err := os.Create(logFilePath)
	if err != nil {
		panic(err)
	}
	writer := zapcore.AddSync(file)
	core := zapcore.NewTee(zapcore.NewCore(fileEncoder, writer, logLevel))
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
