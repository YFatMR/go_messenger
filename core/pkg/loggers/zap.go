package loggers

import (
	"context"
	"os"

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr/logerrcore"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func toZapField(f logerrcore.Field) zap.Field {
	switch f.Type {
	case logerrcore.StringType:
		return zap.String(f.Key, f.String)
	case logerrcore.Int64Type:
		return zap.Int64(f.Key, f.Integer)
	case logerrcore.ErrorType:
		return zap.NamedError(f.Key, f.Error)
	case logerrcore.BoolType:
		val := f.Integer == 1
		return zap.Bool(f.Key, val)
	case logerrcore.SkipType:
	case logerrcore.UnknownType:
	default:
	}
	return zap.Skip()
}

func toZapLogLevel(lvl logerr.Loglevel) zapcore.Level {
	switch lvl {
	case logerr.DebugLevel:
		return zap.DebugLevel
	case logerr.InfoLevel:
		return zap.InfoLevel
	case logerr.ErrorLevel:
	default:
	}
	return zap.ErrorLevel
}

func getZapFormatFields(logerr logerr.Error) []zap.Field {
	if logerr == nil {
		return []zap.Field{}
	}
	result := make([]zap.Field, 0, len(logerr.GetFields()))
	for _, field := range logerr.GetFields() {
		result = append(result, toZapField(field))
	}
	return result
}

func (l *OtelZapLoggerWithTraceID) LogLogerror(lerr logerr.Error) {
	if lerr == nil || !lerr.IsLogMessage() {
		return
	}
	fields := append(getZapFormatFields(lerr), zap.NamedError("api error", lerr.GetAPIError()))
	l.Log(toZapLogLevel(lerr.GetLogLevel()), lerr.GetLogMessage(), fields...)
}

func (l *OtelZapLoggerWithTraceID) LogContextNoExportLogerror(ctx context.Context, lerr logerr.Error) {
	if lerr == nil || !lerr.IsLogMessage() {
		return
	}
	fields := append(getZapFormatFields(lerr), zap.NamedError("api error", lerr.GetAPIError()))
	l.LogContextNoExport(ctx, toZapLogLevel(lerr.GetLogLevel()), lerr.GetLogMessage(), fields...)
}

func (l *OtelZapLoggerWithTraceID) LogContextLogerror(ctx context.Context, lerr logerr.Error) {
	if lerr == nil || !lerr.IsLogMessage() {
		return
	}
	fields := append(getZapFormatFields(lerr), zap.NamedError("api error", lerr.GetAPIError()))
	switch lerr.GetLogLevel() {
	case logerr.DebugLevel:
		l.DebugContext(ctx, lerr.GetLogMessage(), fields...)
	case logerr.InfoLevel:
		l.InfoContext(ctx, lerr.GetLogMessage(), fields...)
	case logerr.ErrorLevel:
		l.ErrorContext(ctx, lerr.GetLogMessage(), fields...)
	}
}
