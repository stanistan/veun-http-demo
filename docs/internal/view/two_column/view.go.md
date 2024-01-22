A two column layout.

We'll have a fixed left navigation and a scrollable
main section.

```go
import (
    _ "embed"
    "context"
    "net/http"

    "github.com/stanistan/veun"
    "github.com/stanistan/veun/vhttp/request"

    "github.com/stanistan/veun-http-demo/internal/view/page"
)
```

## The View

```go
type View struct {
	Nav, Main veun.AsView
	Title     string
	IsMobile  bool
}
```

Embed a template in it instead of using the `el` library,
in case this starts to get more complicated.

```go
//go:embed template.tpl
var tpl string
var template = veun.MustParseTemplate("two_column", tpl)

func (v *View) View(ctx context.Context) (*veun.View, error) {
	return veun.V(veun.Template{
		Tpl:   template,
		Slots: veun.Slots{"nav": v.Nav, "main": v.Main},
		Data:  v,
	}), nil
}
```

### Adding a title

We can have this view provide a title to the page.

```go
func (v *View) SetPageData(d *page.Data) {
    if v.Title != "" {
        d.Title = v.Title
    }

    v.IsMobile = d.IsMobile
}
```

## handler

We can also make our two column view a handler middleware.

```go
type Handler struct {
	Nav, Main request.Handler
}

func (h Handler) ViewForRequest(r *http.Request) (veun.AsView, http.Handler, error) {
	main, next, err := h.Main.ViewForRequest(r)
	if err != nil || main == nil {
		return nil, next, err
	}

	nav, _, err := h.Nav.ViewForRequest(r)
	if err != nil {
		return nil, nil, err
	}

	return &View{Main: main, Nav: nav, Title: "TEST"}, next, nil
}
```
