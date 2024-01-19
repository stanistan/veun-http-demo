```go
import (
    _ "embed"
    "runtime"
    "strings"
    "path/filepath"
    "fmt"
    "context"
    "log/slog"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)
```

## handler

The components handler takes the first name of the components and tries
to find it _both_ as a registered component to render (and how), as well
as the documentation for it by name.

### The component interface

```go
type Component interface {
	veun.AsView
	Title() string
}

type Components []Component

func ComponentView(c Component) component {
	return component{
		Name: fmt.Sprintf("%T - %s", c, c.Title()),
		Body: c,
	}
}

func (v Components) View(ctx context.Context) (*veun.View, error) {
	items := make(veun.Views, len(v))
	for idx, c := range v {
		items[idx] = ComponentView(c)
	}

	return veun.V(items), nil
}

type component struct {
	Name      string
	Body      veun.AsView
	BodyClass string
}

//go:embed component.tpl
var tpl string
var componentTpl = veun.MustParseTemplate("component", tpl)

func (v component) template() veun.Template {
	return veun.Template{
		Tpl:   componentTpl,
		Slots: veun.Slots{"body": v.Body},
		Data:  v,
	}
}

func (v component) View(_ context.Context) (*veun.View, error) {
	// N.B. every single component gets error handling.
	// so no single one gets to break anything else.
	return veun.V(v.template()).WithErrorHandler(v), nil
}

func (v component) ViewForError(ctx context.Context, err error) (veun.AsView, error) {
    slog.Error("error", "err", err)
	return component{
		Name: v.Name,
		Body: veun.Views{
			el.Div().Content(el.Strong().InnerText("Error Captured:")),
            el.P().InnerText(err.Error()),
		},
		BodyClass: "error",
	}, nil
}
```

### component registration

```go
var componentsMap = map[string]Components{}

func c(vs ...Component) {
    _, file, _, ok := runtime.Caller(1)
    if !ok {
        panic("no caller")
    } else {
        file = strings.TrimSuffix(filepath.Base(file), ".generated.go")
    }

	componentsMap[file+".md"] = Components(vs)
}

func ForURL(url string) (veun.AsView, bool) {
	slog.Debug("checking url", "url", url)
    slog.Debug("lets see what we got", "cs", componentsMap)
	v, ok := componentsMap[url]
	return v, ok
}
```
