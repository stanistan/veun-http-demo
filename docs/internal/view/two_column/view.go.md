A two column layout.

We'll have a fixed left navigation and a scrollable
main section.

```go
import (
    _ "embed"
    "context"

    "github.com/stanistan/veun"

    "github.com/stanistan/veun-http-demo/internal/view/page"
)
```

## The View

```go
type View struct {
	Nav, Main veun.AsView
	Title     string
}
```

Embed a template in it instead of using the `el` library,
in case this starts to get more complicated.

```go
//go:embed template.tpl
var tpl string
var template = veun.MustParseTemplate("two_column", tpl)

func (v View) View(ctx context.Context) (*veun.View, error) {
	return veun.V(veun.Template{
		Tpl:   template,
		Slots: veun.Slots{"nav": v.Nav, "main": v.Main},
	}), nil
}
```

### Adding a title

We can have this view provide a title to the page.

```go
func (v View) SetPageData(d *page.Data) {
    d.Title = v.Title
}
```
