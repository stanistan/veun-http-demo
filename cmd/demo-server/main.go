package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/vhttp"
	"github.com/stanistan/veun/vhttp/handler"
	"github.com/stanistan/veun/vhttp/request"

	"github.com/stanistan/veun-http-demo/internal/view"
)

var (
	//go:embed static
	staticFiles embed.FS
	cssPath     = "/static/styles.css"
	htmxPath    = "/static/htmx.1.9.9.min.js"
)

func html(r request.Handler) request.Handler {
	return view.HTML(r, view.HTMLData{
		Title:    "veun-http (demo)",
		CSSPath:  cssPath,
		HTMXPath: htmxPath,
	})
}

var (
	h  = vhttp.Handler
	hf = vhttp.HandlerFunc

	notFound   = h(html(view.NotFoundRequestHandler()))
	fileServer = http.FileServer(http.FS(staticFiles))
)

func main() {
	initLogger()

	mux := http.NewServeMux()

	// our ajax clicked handler that gets hit by htmx
	mux.Handle("/clicked", h(view.ClickTriggerHandler))

	// components and raw components
	mux.Handle("/component/raw", h(view.ComponentPicker()))
	mux.Handle("/component", h(html(view.ComponentPicker())))

	mux.Handle("/components", h(request.Always(view.DefinedComponents)))

	// home view without the html wrapper
	mux.Handle("/home", h(view.HomeViewHandler()))

	// errors
	mux.Handle("/error/banana", hf(func(_ *http.Request) (veun.AsView, http.Handler, error) {
		return nil, nil, fmt.Errorf("banana from /error/banana")
	}))

	mux.Handle("/error/empty", hf(
		// our empty handler just errors out
		func(_ *http.Request) (veun.AsView, http.Handler, error) {
			return nil, nil, fmt.Errorf("banana from /error/empty")
		},
		// our error handler for this one swallows the error
		vhttp.WithErrorHandlerFunc(func(_ context.Context, err error) (veun.AsView, error) {
			slog.Error("error handler called", "err", err)
			return veun.Raw("Custom Error Page"), nil
		}),
	))

	mux.HandleFunc("/r/noop", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, world.\n"))
	})

	mux.HandleFunc("/r/echo", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.URL.Query().Get("in")))
	})

	mux.Handle("/v/noop", h(request.Always(
		veun.Raw("Hello, world.\n"),
	)))

	mux.Handle("/v/echo", hf(func(r *http.Request) (veun.AsView, http.Handler, error) {
		return veun.Raw(r.URL.Query().Get("in")), nil, nil
	}))

	mux.Handle("/lazy", h(html(view.LazyRequestHandler())))

	mux.Handle("/", handler.OnlyRoot(h(html(view.HomeViewHandler()))))

	var addr string
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	} else {
		addr = ":8080"
	}

	s := &http.Server{
		Addr: addr,

		// Our handler does a 404 fallback between the mux, & static files.
		// we introduce a not-found handler as the last http.Handler to check.
		Handler: handler.Checked(
			mux,        // first we we serve our "routes"
			fileServer, // if any of those 404 we check static files
			notFound,   // if any of those 404 we use our notFound handler
		),

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
}

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
