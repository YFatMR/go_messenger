import (
	"github.com/YFatMR/go_messenger/core/pkg/errors/logerr"
)


{{ $errorType := "logerr.Error" }}
{{ $decorator := (or .Vars.DecoratorName (printf "LogerrCleaner%sDecorator" .Interface.Name)) }}

// {{$decorator}} implements {{.Interface.Type}}
// Use {{$decorator}} to make logerr.Error nil if it has no error
// Use {{$decorator}} as last decorator in your chain
type {{$decorator}} struct {
	base {{.Interface.Type}}
}

// New{{$decorator}} instruments an implementation of the {{.Interface.Type}} with simple logging
func New{{$decorator}}(base {{.Interface.Type}}) *{{$decorator}} {
	if base == nil {
		panic("{{$decorator}} got empty base")
	}
	return &{{$decorator}}{
		base: base,
	}
}

{{ range $method := .Interface.Methods }}
	// {{$method.Name}} implements {{$.Interface.Type}}
	func (d *{{$decorator}}) {{$method.Declaration}} {

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
		if {{ $errorResult.Name }} != nil && !{{ $errorResult.Name }}.HasError() {
			{{ $errorResult.Name }} = nil
		}
		{{ $method.Pass "d.base." -}}
	}
{{- end }}
