package main

//go:generate go run github.com/stanistan/veun-http-demo/cmd/lit-gen

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/vhttp"
	"github.com/stanistan/veun/vhttp/handler"
	"github.com/stanistan/veun/vhttp/request"

	"github.com/stanistan/veun-http-demo/internal/view"
)

func serverOld() http.Handler {

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

	mux.Handle("/docs", h(html(docsIndex)))
	mux.Handle("/docs/", http.StripPrefix("/docs/", h(html(docsPage))))

	// mux.Handle("/", handler.OnlyRoot(h(html(view.HomeViewHandler()))))
	mux.Handle("/", handler.OnlyRoot(h(html(index))))

	// Our handler does a 404 fallback between the mux, & static files.
	// we introduce a not-found handler as the last http.Handler to check.
	return handler.Checked(
		mux,                                    // first we we serve our "routes"
		staticFileServer(),                     // if any of those 404 we check static files
		h(html(view.NotFoundRequestHandler())), // if any of those 404 we use our notFound handler
	)
}
