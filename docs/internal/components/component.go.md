We build a standard component to demo some things
that you can do with things the library.

```go
import (
	"context"
	_ "embed"
	"fmt"

	"github.com/stanistan/veun"
	t "github.com/stanistan/veun/template"
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
		Type:        fmt.Sprintf("%T", c),
		Description: c.Description(),
		Body:        c,
	}
}
```

Our component uses the `component.tpl` in this directory.

```go
//go:embed component.tpl
var tpl string
var componentTpl = t.MustParse("component", tpl)

type component struct {
	Type, Description string

	Body      veun.AsView
    BodyClass string
}
```

And the view that we create has an error handler based on the
component itself.

```go
func (v component) View(ctx context.Context) (*veun.View, error) {
	return veun.V(v.template()).WithErrorHandler(errorHandler{v}), nil
}
```

The `template` is reusable so that we can embed it into the
[`errorView`](/docs/internal/components/error.md).

```go
func (v component) template() t.Template {
	return t.Template{
		Tpl:   componentTpl,
		Slots: t.Slots{"body": v.Body},
		Data:  v,
	}
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
