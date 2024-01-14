# routes

```go
import (
	"net/http"

	"github.com/stanistan/veun/vhttp"
	"github.com/stanistan/veun/vhttp/handler"

	"github.com/stanistan/veun-http-demo/internal/view/page"
)
```

## Helpers

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
    Title:    "veun-http (demo)",
    CSSFiles: []string{"/static/styles.css"},
    JSFiles:  []string{"/static/htmx.1.9.9.min.js", "/static/prism.js"},
})
```

Any route can overwrite and/or add data to this struct if it
implements `page.DataMutator`.

## routes / handler

Open our function body:

```go
func routes() http.Handler {
    mux := http.NewServeMux()
```

---

### Docs!

We only serve two routes for docs, and they can use the same
handler for both.

```go
mux.Handle("/docs", h(html(docsHandler)))
mux.Handle("/docs/", h(html(docsHandler)))
```

I really feel like there should be a better pattern for this when
using `http.ServeMux`, but for now, the repetition is _fine_.

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
return handler.Checked(mux, staticFileServer())
```

---

Closing `server`:

```go
}
```
