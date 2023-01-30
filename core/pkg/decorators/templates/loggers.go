import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/cerrors"
	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"go.uber.org/zap"
)

{{ $decorator := (or .Vars.DecoratorName (printf "Logging%sDecorator" .Interface.Name)) }}

// {{$decorator}} implements {{.Interface.Type}} that is instrumented with custom zap logger
type {{$decorator}} struct {
	logger *loggers.OtelZapLoggerWithTraceID
	base {{.Interface.Type}}
}

// New{{$decorator}} instruments an implementation of the {{.Interface.Type}} with simple logging
func New{{$decorator}}(base {{.Interface.Type}}, logger *loggers.OtelZapLoggerWithTraceID) *{{$decorator}} {
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
		panic("Expected return type from method. cerrors.Error at least")
		{{ break }}
	{{ end }}

	{{ $errorsResultCount := 0}}
	{{ range $result := $method.Results }}
		{{ if eq $result.Type "cerrors.Error" }}
			{{ $errorsResultCount = add $errorsResultCount 1 }}
		{{ end }}
	{{ end }}

	{{ if ne $errorsResultCount 1 }}
		panic("Expected exact one cerrors.Error type as last argument")
		{{ break }}
	{{ end }}

	{{ $errorResult := last $method.Results }}
	{{ if not (eq $errorResult.Type "cerrors.Error") }}
		panic("Expected exact one cerrors.Error type as last argument")
		{{ break }}
	{{ end }}
		d.logger.DebugContextNoExport(ctx, "{{$decorator}}: calling {{$method.Name}}")
		defer func() {
			if {{ $errorResult.Name }} != nil {
				d.logger.ErrorContext(
					ctx, {{ $errorResult.Name }}.GetInternalErrorMessage(), zap.Error({{ $errorResult.Name }}.GetInternalError()),
				)
			} else {
				d.logger.DebugContextNoExport(ctx, "{{$decorator}}: {{$method.Name}} finished")
			}
			d.logger.DebugContextNoExport(ctx, "{{$decorator}}: {{$method.Name}} finished")
		}()
		{{ $method.Pass "d.base." -}}
	}
{{- end }}
