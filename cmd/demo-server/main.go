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
	orNotFound = func(next http.Handler) http.Handler {
		return handler.Checked(next, notFound)
	}
)

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	)))

	// 1.
	mux := http.NewServeMux()

	// our ajax clicked handler that gets hit by htmx
	mux.Handle("/clicked", h(view.ClickTriggerHandler))

	// components and raw components
	mux.Handle("/component/raw", h(view.ComponentPicker(notFound)))
	mux.Handle("/component", h(html(view.ComponentPicker(notFound))))

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

	mux.Handle("/", handler.Checked(
		// 1. First we check our root page
		handler.OnlyRoot(h(html(view.HomeViewHandler()))),

		// 2. Then we see if we can find things in our static files
		http.FileServer(http.FS(staticFiles)),

		// 3. What to do with our NotFound
		h(html(view.NotFoundRequestHandler())),
	))

	var addr string
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	} else {
		addr = ":8080"
	}

	s := &http.Server{
		Addr:    addr,
		Handler: mux,

		// logging
		ErrorLog: slog.NewLogLogger(slog.Default().Handler(), slog.LevelWarn),

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
