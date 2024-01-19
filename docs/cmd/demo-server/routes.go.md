```go
import (
	"net/http"

	"github.com/stanistan/veun/el"
	"github.com/stanistan/veun/vhttp"
	"github.com/stanistan/veun/vhttp/handler"
	"github.com/stanistan/veun/vhttp/request"

	"github.com/stanistan/veun-http-demo/internal/view/page"
)
```

## helpers

Becuase we're going to be doing lots of composition of
request handers and http handlers (from `veun`), we have some
helpers defined.

These two are to wrap a `request.Handler(Func)` and create an
`http.Handler`.

```go
var (
	h  = vhttp.Handler
	hf = vhttp.HandlerFunc
)
```

### html

This is a wrapper for our html page, with default title,
css, and js files provided. `page.Handler` returns
a middleware-y function.

```go
var html = page.Handler(page.Data{
	Title:    "veun-http-demo", // default title
	CSSFiles: []string{"/static/styles.css"},
	JSFiles:  []string{"/static/htmx.1.9.9.min.js", "/static/prism.js"},
})
```

Any view/request.Handler can mutate the data we pass here by
implementing the `page.DataMutator` interface.

## handlers

Open our function body:

```go
func routes() http.Handler {
    mux := http.NewServeMux()
```

### Docs & Components

Docs are really the only interesting route here. It serves everything,
from the static documents to demo compontents in the `internal/components`
path.

```go
mux.Handle("/docs/", h(html(docsHandler)))
```

#### Pages for lazy loading

```go
mux.Handle("/e/text", h(request.Always(el.Div().Content(
    el.Em().InnerText("text output,").In(el.P()),
    el.P().InnerText("and more of it"),
))))
```

```go
mux.Handle("/e/not_found", h(request.Always(notFoundView)))
```

### Root

`handler.OnlyRoot` is in the `veun` library that ensures
that when we're mounting at the `/` path, we 404 if it's _anything_
other than that exact path.

```go
mux.Handle("/", handler.OnlyRoot(h(html(index))))
```

### Closing out the server

There is another built in handler `handler.Checked`
which will continue trying `http.Handler`s if it hits a 404 and
return the last one.

This way, we can implement a pretty nice static file serving fallback for
when any routes aren't actually defined, and even add our own 404
handler, it's pretty neat.

```go
return handler.Checked(
	mux,                // the ServeMux we've just been adding routes to
	staticFileServer(), // falls back to the static file server if we 404
	notFoundHandler(),  // our own custom notFoundHandler
)
```

Closing `server`:

```go
}
```

