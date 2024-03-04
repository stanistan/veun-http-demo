```go
import (
    "path/filepath"

    "github.com/stanistan/veun/el-exp"
)
```

I want to have a bread-crumby title view that's friendlier to
mobile clients.

```go
func View(urlPath string) el.Div {
	dir, file := filepath.Split(urlPath)
    return el.Div{
        el.Class("page-title"),
        el.H1{
            el.Span{el.Text(file)},
            el.Span{
                el.Class("sub-title"),
                el.Text("in: " + dir),
            },
        },
    }
}
```
