import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"go.uber.org/zap"
)

{{ $logObjectType := "ulo.LogStash" }}
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

	{{ if not $method.ReturnsError }}
		panic("Expected return type error")
		{{ break }}
	{{ end }}

	{{ $logObjectCount := 0}}
	{{ $logObjectName := "" }}
	{{ range $result := $method.Results }}
		{{ if eq $result.Type $logObjectType }}
			{{ $logObjectCount = add $logObjectCount 1 }}
			{{ $logObjectName = $result.Name }}
		{{ end }}
	{{ end }}

	{{ if ne $logObjectCount 1 }}
		panic("Expected exact one instance of {{ $logObjectType }} type")
		{{ break }}
	{{ end }}

		d.logger.DebugContextNoExport(ctx, "{{ $decorator }}: calling {{ $method.Name }}")
		defer func() {
			d.logger.LogContextNoExportULO(ctx, {{ $logObjectName }})
			if err != nil {
				d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
			}
			d.logger.DebugContextNoExport(ctx, "{{$decorator}}: {{$method.Name}} finished")
		}()
		{{ $method.Pass "d.base." -}}
	}
{{- end }}
