package loggers

import (
	"context"
	. "core/pkg/utils"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
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

func NewBaseOtelZapLogger(logLevel zapcore.LevelEnabler, logFilePath string) (*otelzap.Logger, error) {
	logger, err := NewBaseZapFileLogger(logLevel, logFilePath)
	if err != nil {
		return nil, err
	}
	return otelzap.New(
		logger,
		otelzap.WithTraceIDField(true),
		otelzap.WithMinLevel(zapcore.ErrorLevel),
		otelzap.WithStackTrace(true),
	), nil
}

// OtelZapLoggerWithTraceID expands the capabilities of the logger otelzap.
// otelzap can only add trace_id to messages that will be passed to the exporter.
// Functions with suffix `...NoExport` resolve this problem
type OtelZapLoggerWithTraceID struct {
	*otelzap.Logger
}

func NewOtelZapLoggerWithTraceID(logLevel zapcore.LevelEnabler, logFilePath string) (*OtelZapLoggerWithTraceID, error) {
	logger, err := NewBaseOtelZapLogger(logLevel, logFilePath)
	if err != nil {
		return nil, err
	}
	return &OtelZapLoggerWithTraceID{
		Logger: logger,
	}, nil
}

func (l *OtelZapLoggerWithTraceID) LogContextNoExport(ctx context.Context, level zapcore.Level, msg string, fields ...zapcore.Field) {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()
	fields = append(fields, zap.String("trace_id", traceID))
	l.Log(level, msg, fields...)
}

// DebugContextNoExport Provide an ability to write trace ID without exporting
func (l *OtelZapLoggerWithTraceID) DebugContextNoExport(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.LogContextNoExport(ctx, zapcore.DebugLevel, msg, fields...)
}

// InfoContextNoExport Provide an ability to write trace ID without exporting
func (l *OtelZapLoggerWithTraceID) InfoContextNoExport(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.LogContextNoExport(ctx, zapcore.InfoLevel, msg, fields...)
}

// WarningContextNoExport Provide an ability to write trace ID without exporting
func (l *OtelZapLoggerWithTraceID) WarningContextNoExport(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.LogContextNoExport(ctx, zapcore.WarnLevel, msg, fields...)
}

// ErrorContextNoExport Provide an ability to write trace ID without exporting
func (l *OtelZapLoggerWithTraceID) ErrorContextNoExport(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.LogContextNoExport(ctx, zapcore.ErrorLevel, msg, fields...)
}

// FatalContextNoExport Provide an ability to write trace ID without exporting
func (l *OtelZapLoggerWithTraceID) FatalContextNoExport(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.LogContextNoExport(ctx, zapcore.FatalLevel, msg, fields...)
}
