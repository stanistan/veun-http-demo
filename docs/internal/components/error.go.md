```go
import (
	"context"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)
```

### handler

This is our error handler for components. The `errorHandler` is parameterized
by a component (it's specific for one) so that we can have a _nice_ looking
error view.

```go
type errorHandler struct {
    c component
}

func (e errorHandler) ViewForError(ctx context.Context, err error) (veun.AsView, error) {
	return errorView{
		c: component{
			// 1. Replace the body of the component with an error
			Body: errorBody("Error Captured by component:", err.Error()),
			// 2. Add a css class
			BodyClass: "error",
			// 3. Change the type to indicate an error as well
			Type: e.c.Type + " !!FAILED!!",
            // 4. Keep the rest...
            Description: e.c.Description,
		},
	}, nil
}
```

### view

The view re-uses the structure of the component struct. The big difference
is that it doesn't also have an error handler. Having a recursive error
handler could cause a big problem.

```go
type errorView struct {
    c component
}

func (v errorView) View(ctx context.Context) (*veun.View, error) {
    return veun.V(v.c.template()), nil
}
```

And our error body:

```go
func errorBody(title, content string) veun.AsView {
    return el.Div{
        el.Strong{el.Text(title)},
        el.P{el.Text(content)},
    }
}
```
