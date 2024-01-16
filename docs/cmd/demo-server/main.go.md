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
this is what we defined in the [boostrapping doc](/docs/cmd/demo-server/bootstrap).

```go
initLogger()
```

The server error logger is also configured to use the `slog` logger.

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

You'll notice the `routes` definition in the call,
this is defined [here](/docs/cmd/demo-server/routes).

This is a function that returns a standard `http.Handler`.

This `main` function body is just boilerpate for serving requests.

```go
s := &http.Server{
	Addr:    addr,
	Handler: routes(), // <- the handler!

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
