package loggers

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Logger initialization

func NewJsonEncoderLogger(logLevel zapcore.LevelEnabler, logFilePath string) *zap.Logger {
	// syncer
	file, err := os.Create(logFilePath)
	if err != nil {
		panic(err)
	}
	writerSyncer := zapcore.AddSync(file)

	// encoder
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	core := zapcore.NewCore(encoder, writerSyncer, logLevel)
	return zap.New(core)
}
