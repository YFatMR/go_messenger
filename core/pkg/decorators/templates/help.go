import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"go.uber.org/zap"
)

{{ $errorType := "logerr.Error" }}
{{ $decorator := (or .Vars.DecoratorName (printf "Logging%sDecorator" .Interface.Name)) }}

// {{$decorator}} implements {{.Interface.Type}} that is instrumented with custom zap logger
type {{$decorator}} struct {
	logger *loggers.OtelZapLoggerWithTraceID
	base {{.Interface.Type}}
}

// New{{$decorator}} instruments an implementation of the {{.Interface.Type}} with simple logging
func New{{$decorator}}(base {{.Interface.Type}}, logger *loggers.OtelZapLoggerWithTraceID) *{{$decorator}} {
	if base == nil {
		panic("{{$decorator}} got empty base")
	}
	if logger == nil {
		panic("{{$decorator}} got empty logger")
	}
	return &{{$decorator}}{
		base: base,
		logger: logger,
	}
}

{{ range $method := .Interface.Methods }}
	// {{$method.Name}} implements {{$.Interface.Type}}
	func (d *{{$decorator}}) {{$method.Declaration}} {

	{{ if not ($method.AcceptsContext) }}
		panic("Expected context variable")
		{{ break }}
	{{ end }}
	{{ if not ($method.HasResults) }}
		panic("Expected return type from method. {{ $errorType }} at least")
		{{ break }}
	{{ end }}

	{{ $errorsResultCount := 0}}
	{{ range $result := $method.Results }}
		{{ if eq $result.Type $errorType }}
			{{ $errorsResultCount = add $errorsResultCount 1 }}
		{{ end }}
	{{ end }}

	{{ if ne $errorsResultCount 1 }}
		panic("Expected exact one {{ $errorType }} type as last argument")
		{{ break }}
	{{ end }}

	{{ $errorResult := last $method.Results }}
	{{ if not (eq $errorResult.Type $errorType) }}
		panic("Expected exact one {{ $errorType }} type as last argument")
		{{ break }}
	{{ end }}
		d.logger.DebugContextNoExport(ctx, "{{ $decorator }}: calling {{ $method.Name }}")
		defer func() {
			defer d.logger.DebugContextNoExport(ctx, "{{$decorator}}: {{$method.Name}} finished")
			if {{ $errorResult.Name }} == nil {
				return
			}
			if {{ $errorResult.Name }}.IsLogMessage() {
				d.logger.LogContextLogerror(ctx, {{ $errorResult.Name }})
				{{ $errorResult.Name }}.StopLogMessage()
			}
		}()
		{{ $method.Pass "d.base." -}}
	}
{{- end }}
