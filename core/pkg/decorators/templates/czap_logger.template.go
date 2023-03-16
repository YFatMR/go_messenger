import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"go.uber.org/zap"
)

{{ $logObjectType := "ulo.LogStash" }}
{{ $decorator := (or .Vars.DecoratorName (printf "Logging%sDecorator" .Interface.Name)) }}

// {{$decorator}} implements {{.Interface.Type}} that is instrumented with custom zap logger
type {{$decorator}} struct {
	logger *czap.Logger
	base {{.Interface.Type}}
}

// New{{$decorator}} instruments an implementation of the {{.Interface.Type}} with simple logging
func New{{$decorator}}(base {{.Interface.Type}}, logger *czap.Logger) *{{$decorator}} {
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
		d.logger.Info("{{ $decorator }}: calling {{ $method.Name }}")
		defer d.logger.Info("{{$decorator}}: {{$method.Name}} finished")
		{{ $method.Pass "d.base." -}}
	}
		{{ continue }}
	{{ end }}

	{{ if not $method.ReturnsError }}
		panic("Expected return type error")
		{{ break }}
	{{ end }}

		d.logger.InfoContext(ctx, "{{ $decorator }}: calling {{ $method.Name }}")
		defer func() {
			if err != nil {
				d.logger.ErrorContext(ctx, "", zap.NamedError("public api error", err))
			}
			d.logger.InfoContext(ctx, "{{$decorator}}: {{$method.Name}} finished")
		}()
		{{ $method.Pass "d.base." -}}
	}
{{- end }}
