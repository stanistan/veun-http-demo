```go
import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/stanistan/veun"
)
```

_We need a way to embed the components in the actual page._

A static map and an init hook is enough for this, and our
doc handler can use this to see if there's any component
to actually render in the page.


```go
var registry = map[string]Views{}

func show(vs ...IComponent) {
    _, file, _, ok := runtime.Caller(1)
    if !ok {
        panic("no caller")
    } else {
        file = strings.TrimSuffix(filepath.Base(file), ".generated.go")
    }

	registry[file+".md"] = Views(vs)
}
```

## Access

Callers should use `ForURL` to get a view based on the key.

```go
func ForURL(url string) (veun.AsView, bool) {
	v, ok := registry[url]
	return v, ok
}
```
