package czap

import (
	"context"
	"os"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const TraceIDKey = "trace_id"

type Settings struct {
	LogTraceID            bool
	ExportMessageLogLevel zapcore.Level
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
type Logger struct {
	otelzap.Logger
	settings Settings
}

func New(logger otelzap.Logger, settings Settings) *Logger {
	return &Logger{
		Logger:   logger,
		settings: settings,
	}
}

func NewNop() *Logger {
	return &Logger{
		Logger: *otelzap.New(zap.NewNop()),
	}
}

func FromConfig(config *cviper.CustomViper) (*Logger, error) {
	logLevel := config.GetZapcoreLogLevelRequired("LOG_LEVEL")
	logPath := config.GetStringRequired("LOG_PATH")

	// Init logger
	zapLogger, err := NewBaseZapFileLogger(logLevel, logPath)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: *otelzap.New(
			zapLogger,
			otelzap.WithTraceIDField(true),
			otelzap.WithMinLevel(zapcore.ErrorLevel),
			otelzap.WithStackTrace(true),
		),
		settings: Settings{
			LogTraceID:            true,
			ExportMessageLogLevel: zap.ErrorLevel,
		},
	}, nil
}

func (l *Logger) logWithContext(ctx context.Context, level zapcore.Level,
	msg string, fields ...zapcore.Field,
) {
	if l.settings.LogTraceID {
		span := trace.SpanFromContext(ctx)
		traceID := span.SpanContext().TraceID().String()
		fields = append(fields, zap.String(TraceIDKey, traceID))
	}
	l.Log(level, msg, fields...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	exportLogMessage := l.settings.ExportMessageLogLevel <= zap.DebugLevel
	if l.settings.LogTraceID && exportLogMessage {
		l.Logger.DebugContext(ctx, msg, fields...)
		return
	}
	l.logWithContext(ctx, zapcore.DebugLevel, msg, fields...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	exportLogMessage := l.settings.ExportMessageLogLevel <= zap.InfoLevel
	if l.settings.LogTraceID && exportLogMessage {
		l.Logger.InfoContext(ctx, msg, fields...)
		return
	}
	l.logWithContext(ctx, zapcore.InfoLevel, msg, fields...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	exportLogMessage := l.settings.ExportMessageLogLevel <= zap.WarnLevel
	if l.settings.LogTraceID && exportLogMessage {
		l.Logger.WarnContext(ctx, msg, fields...)
		return
	}
	l.logWithContext(ctx, zapcore.WarnLevel, msg, fields...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	exportLogMessage := l.settings.ExportMessageLogLevel <= zap.ErrorLevel
	if l.settings.LogTraceID && exportLogMessage {
		l.Logger.ErrorContext(ctx, msg, fields...)
		return
	}
	l.logWithContext(ctx, zapcore.ErrorLevel, msg, fields...)
}

func (l *Logger) DPanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	exportLogMessage := l.settings.ExportMessageLogLevel <= zap.DPanicLevel
	if l.settings.LogTraceID && exportLogMessage {
		l.Logger.DPanicContext(ctx, msg, fields...)
		return
	}
	l.logWithContext(ctx, zapcore.DPanicLevel, msg, fields...)
}

func (l *Logger) PanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	exportLogMessage := l.settings.ExportMessageLogLevel <= zap.PanicLevel
	if l.settings.LogTraceID && exportLogMessage {
		l.Logger.PanicContext(ctx, msg, fields...)
		return
	}
	l.logWithContext(ctx, zapcore.PanicLevel, msg, fields...)
}

func (l *Logger) FatalContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	exportLogMessage := l.settings.ExportMessageLogLevel <= zap.FatalLevel
	if l.settings.LogTraceID && exportLogMessage {
		l.Logger.FatalContext(ctx, msg, fields...)
		return
	}
	l.logWithContext(ctx, zapcore.FatalLevel, msg, fields...)
}
