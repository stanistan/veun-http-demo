# html page

I want some basic styling and syntax highlighting for
my webpage, and that means we actually need to include that in
some standard "chrome" wrapper for everything else.

## Imports, obv

```go
import (
    "github.com/stanistan/veun/vhttp/request"

    "github.com/stanistan/veun-http-demo/internal/view/page"
)
```

## Where our static files live

```go
const (
	cssPath     = "/static/styles.css"
	htmxPath    = "/static/htmx.1.9.9.min.js"
	prismJSPath = "/static/prism.js"
)
```

## Composition of htmlPage

```go
func htmlPage(r request.Handler) request.Handler {
	return page.Handler(r, page.Data{
		Title:    "veun-http (demo)",
		CSSFiles: []string{cssPath},
		JSFiles:  []string{htmxPath, prismJSPath},
	})
}
```
