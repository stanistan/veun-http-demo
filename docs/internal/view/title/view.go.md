```go
import (
    "path/filepath"

    "github.com/stanistan/veun/el"
)
```

I want to have a bread-crumby title view that's friendlier to
mobile clients.

```go
func View(urlPath string) *el.Element {
	dir, file := filepath.Split(urlPath)
	return el.Div().Class("page-title").Content(
		el.H1().Content(
			el.Span().InnerText(file),
			el.Span().Class("sub-title").InnerText("in: "+dir),
		),
	)
}
```
