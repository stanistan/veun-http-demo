# main

We're going to set up our main function, or at least
stub it out and then define our server definitions.

```go
import (
    "log/slog"
    "net/http"
    "os"
    "time"
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
	Handler: server(), // server! we'll fill this in the next one

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
