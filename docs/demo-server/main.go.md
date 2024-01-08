# main

We're going to set up our main function, or at least
stub it out and then define our server definitions.

```go
import (
    "log/slog"
    "net/http"
    "os"
    "time"

    "github.com/stanistan/veun/vhttp"
    "github.com/stanistan/veun/vhttp/handler"
)

func main() {
```

## The Deets

Let's actually call our logger first, this part is uninteresting,
this is what we defined in the [first doc](/doc/1).

```go
initLogger()
```

### Server Address

And we can get the PORT we're running on from the environment,
defaulting to `8080`.

```go
var addr string
if port := os.Getenv("PORT"); port != "" {
	addr = ":" + port
} else {
	addr = ":8080"
}
```


### Defining the server

You'll notice the stub `server` definition in the call,
for now, we can think of it as an empty `func() http.Handler`
but we'll end up adding things to it (like our routes) so that
this whole server actually does _something_.

For now this has all definitely been boilerplate to get us set
up with static file serving, and just regular serving, but
we're not _using_ any of that just yet.

```go
s := &http.Server{
	Addr:    addr,
	Handler: server(),

	// logging
	ErrorLog: slog.NewLogLogger(
		slog.Default().Handler(), slog.LevelWarn,
	),

	// timeouts
	ReadTimeout:  1 * time.Second,
	WriteTimeout: 5 * time.Second,
	IdleTimeout:  5 * time.Second,
}

slog.Info("starting server", "addr", addr)
if err := s.ListenAndServe(); err != nil {
	slog.Error("server stopped", slog.String("err", err.Error()))
}
```

---

Closing `main`:

```go
}
```

## Defining our routes

### Helpers

Becuase we're going to be doing lots of composition of
request handers and http handlers (from `veun`), we have some
helpers defined.

```go
var (
	h  = vhttp.Handler
	hf = vhttp.HandlerFunc
)

func server() http.Handler {
    mux := http.NewServeMux()
```

### Docs routes

```go
mux.Handle("/docs", h(htmlPage(docsIndex)))
mux.Handle("/docs/", http.StripPrefix("/docs/", h(htmlPage(docsPage))))
```

### Root

`handler.OnlyRoot` is in the `veun` library that ensures
that when we're mounting at the `/` path, we 404 if it's _anything_
other than that exact path.

```go
mux.Handle("/", handler.OnlyRoot(h(htmlPage(index))))
```

And we have another built in handler `Checked`
which will continue trying `http.Handler`s if it hits a 404.

This way, we can implement a pretty nice static file serving fallback for
when any routes aren't actually defined.

```go
return handler.Checked(mux, staticFileServer())
```

### Backing off on 404s to

---

Closing `server`:

```go
}
```
