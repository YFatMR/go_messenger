package loggers

import (
	"context"
	"os"

	"github.com/YFatMR/go_messenger/core/pkg/utils"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func RequiredZapcoreLogLevelEnv(env string) zapcore.Level {
	level := utils.RequiredStringEnv(env)
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		panic("Variable " + env + " has unexpected values")
	}
}

func NewBaseZapFileLogger(logLevel zapcore.LevelEnabler, logFilePath string) (*zap.Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)

	file, err := os.Create(logFilePath)
	if err != nil {
		return nil, err
	}
	writer := zapcore.AddSync(file)
	core := zapcore.NewTee(zapcore.NewCore(fileEncoder, writer, logLevel))
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}

// OtelZapLoggerWithTraceID expands the capabilities of the logger otelzap.
// otelzap can only add trace_id to messages that will be passed to the exporter.
// Functions with suffix `...NoExport` resolve this problem.
type OtelZapLoggerWithTraceID struct {
	*otelzap.Logger
}

func NewOtelZapLoggerWithTraceID(logger *otelzap.Logger) *OtelZapLoggerWithTraceID {
	return &OtelZapLoggerWithTraceID{
		Logger: logger,
	}
}

// LogContextNoExport Provide an ability to write trace ID without exporting.
func (l *OtelZapLoggerWithTraceID) LogContextNoExport(ctx context.Context, level zapcore.Level,
	msg string, fields ...zapcore.Field,
) {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()
	fields = append(fields, zap.String("trace_id", traceID))
	l.Log(level, msg, fields...)
}

// DebugContextNoExport Provide an ability to write trace ID without exporting.
func (l *OtelZapLoggerWithTraceID) DebugContextNoExport(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.LogContextNoExport(ctx, zapcore.DebugLevel, msg, fields...)
}

// InfoContextNoExport Provide an ability to write trace ID without exporting.
func (l *OtelZapLoggerWithTraceID) InfoContextNoExport(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.LogContextNoExport(ctx, zapcore.InfoLevel, msg, fields...)
}

// WarningContextNoExport Provide an ability to write trace ID without exporting.
func (l *OtelZapLoggerWithTraceID) WarningContextNoExport(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.LogContextNoExport(ctx, zapcore.WarnLevel, msg, fields...)
}

// ErrorContextNoExport Provide an ability to write trace ID without exporting.
func (l *OtelZapLoggerWithTraceID) ErrorContextNoExport(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.LogContextNoExport(ctx, zapcore.ErrorLevel, msg, fields...)
}

// FatalContextNoExport Provide an ability to write trace ID without exporting.
func (l *OtelZapLoggerWithTraceID) FatalContextNoExport(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.LogContextNoExport(ctx, zapcore.FatalLevel, msg, fields...)
}
