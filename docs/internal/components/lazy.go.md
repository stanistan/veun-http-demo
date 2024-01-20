```go
import (
	"context"
	"errors"
	"fmt"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)
```

## Lazy Loading

When building websites, sometimes it's really nice
to have the chunk of the page load _lazilly_ and _separately_
from everything else.

We have [htmx][htmx] loaded client side, so we can leverage that
and plain old http/html responses to accomplish this.

Our struct takes a URL to load content from, a placeholder view,
and a delay.

```go
type Lazy struct {
	URL         string
	Placeholder veun.AsView

	Delay string
}
```

How does this render?

```go
func (v Lazy) View(ctx context.Context) (*veun.View, error) {
    if v.URL == "" {
        return nil, errors.New("no url")
    }

    return el.Div().
        Attrs(v.htmxAttrs()).
        Content(v.placeholder()).
        View(ctx)
}
```

We get our attributes for `htmx`:

```go
func (v Lazy) htmxAttrs() el.Attrs {
	trigger := "load"
	if v.Delay != "" {
		trigger += " delay:" + v.Delay
	}
	return el.Attrs{"hx-get": v.URL, "hx-trigger": trigger}
}
```

And a default placeholder:

```go
func (v Lazy) placeholder() veun.AsView {
    if v.Placeholder == nil {
        return el.Em().InnerText("...loading...")
    } else {
        return v.Placeholder
    }
}
```

### A component

`Lazy` is a component...

```go
var _ IComponent = Lazy{}

func (v Lazy) Description() string {
    return fmt.Sprintf("url=%s delay=%s placeholder=%+v", v.URL, v.Delay, v.Placeholder)
}
```

## registering it

```go
func init() {
	show(
		Lazy{URL: "/not_found"},
		Lazy{URL: "/e/text"},
		Lazy{URL: "/e/text", Placeholder: el.Text("custom placeholder... with delay"), Delay: "5s"},
		Lazy{}, // Missing URL
	)
}
```

[htmx]: https://htmx.org/
