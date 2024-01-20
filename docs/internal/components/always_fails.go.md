The components you see above are rendered by the `init`
register hook and are declared below. This demonstrates
how you can do error handling with `veun`.

To do this, we create a view that always fails. It's configured
to either capture its own error, or not. But it always fails.

When it doesn't capture its own error, the rendering parent
component captures it. If we left it alone, we'd end up with a
500 server error.

```go
import (
	"context"
	"errors"

	"github.com/stanistan/veun"
)
```

## init component registration

What renders the components at the top of this page :)

```go
func init() {
    show(AlwaysFails{OwnErrorCapture: true}, AlwaysFails{OwnErrorCapture: false})
}
```

## the view

```go
type AlwaysFails struct {
	OwnErrorCapture bool
}
```

`AlwaysFails` has two implementations based on `OwnErrorCapture`.

```go
func (v AlwaysFails) View(_ context.Context) (*veun.View, error) {
```

Either we propagate an error out of the view directly:

```go
if !v.OwnErrorCapture {
    return nil, errors.New("this view will always fail")
}
```

Or we create uh kind of recursively create view that captures
the error of _itself_.

```go
// N.B. Why yes, this is a recursive definition.
return veun.V(AlwaysFails{}).WithErrorHandler(v), nil
```

Closing `View`:

```go
}
```

### Error Handler

As we see above, this view is also an error handler.

```go
var _ veun.ErrorHandler = AlwaysFails{}

func (v AlwaysFails) ViewForError(_ context.Context, err error) (veun.AsView, error) {
    return errorBody("Error, captured by AlwaysFails:", err.Error()), nil
}
```

### Components

We fulfill the `Component` interface by giving this a description in the component UI.

```go
func (v AlwaysFails) Description() string {
	if v.OwnErrorCapture {
		return "captures itself"
	}

	return "captured by component"
}
```

