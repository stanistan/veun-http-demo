We build a standard component to demo some things
that you can do with things the library.

```go
import (
	"context"
	_ "embed"
	"fmt"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)
```

I don't know exactly why I feel like prefacing the name of the
interface with `I` here, but given the overloading
of _component_, I'm going to do it.

Each component is first and foremost a `veun.AsView`, and we attach
behavior to it. For us the behavior is silly, we get metadata, a
description.

```go
type IComponent interface {
	veun.AsView
	Description() string
}
```

We also have a _concrete_ view that we're creating from the interface.

```go
func View(c IComponent) component {
    return component{
        Type: fmt.Sprintf("%T", c),
        Description: c.Description(),
        Body: c,
    }
}

type component struct {
	Type, Description string

	Body      veun.AsView
	BodyClass string
}

//go:embed component.tpl
var tpl string
var componentTpl = veun.MustParseTemplate("component", tpl)

func (v component) View(ctx context.Context) (*veun.View, error) {
	return veun.V(veun.Template{
		Tpl:   componentTpl,
		Slots: veun.Slots{"body": v.Body},
		Data:  v,
	}).WithErrorHandler(v), nil
}
```

A component is its own error handler and can render a component
with an error, so any failure will not break the page/request
context in which a "component" is used.

Note the addition of the `error` body class which impacts the template.

```go
func (v component) ViewForError(ctx context.Context, err error) (veun.AsView, error) {
    // FIXME: having something be an error handler for itself is a bad idea
    // you can get into a recursive error thing and maybe the library prevents
    // this?
	return component{
		Type:        v.Type,
		Description: "err: " + v.Description,
		Body: veun.Views{
			el.Div().Content(el.Strong().InnerText("Error Captured:")),
			el.P().InnerText(err.Error()),
		},
		BodyClass: "error",
	}, nil
}
```

### many components

We can create our own _components_ container, similar to `veun.Views`.

```go
type Views []IComponent

func (v Views) View(ctx context.Context) (*veun.View, error) {
    items := make(veun.Views, len(v))
    for idx, c := range v {
        items[idx] = View(c)
    }

    return items.View(ctx)
}
```
