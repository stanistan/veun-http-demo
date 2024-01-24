# Bootstrapping

The first couple of parts/docs here are to set up a scaffolding
for the work we're going to demo using `veun`.

## stdlib

Let's start with building out the imports we're going to need for our
main sever implementation from the standard library.

```go
import (
	"embed"
	"log/slog"
	"net/http"
	"os"
)
```

## Static Files

We're going to embed the entire static directory into
this server binary and define the css path and htmx path
to add that into our served html page as we go.

```go
var (
	//go:embed static
	staticFiles embed.FS
)
```

And we're going to add a static fileserver handler
based on `net/http` as well.

```go
func staticFileServer() http.Handler {
	return http.FileServer(http.FS(staticFiles))
}
```

## Logging

What do we do with the logger? Set up a structured
logging based on an `ENV` environment variable.

```go
func initLogger() {
	var logHandler slog.Handler
	if os.Getenv("ENV") == "dev" {
		logHandler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	slog.SetDefault(slog.New(logHandler))
}
```
