package main

//go:generate go run github.com/stanistan/veun-http-demo/cmd/lit-gen -root ../../docs/demo-server -o .

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/html"
	"github.com/stanistan/veun/vhttp"
	"github.com/stanistan/veun/vhttp/handler"
	"github.com/stanistan/veun/vhttp/request"

	"github.com/stanistan/veun-http-demo/docs"
	"github.com/stanistan/veun-http-demo/internal/view"
	"github.com/stanistan/veun-http-demo/internal/view/md"
	"github.com/stanistan/veun-http-demo/internal/view/page"
)

const (
	cssPath     = "/static/styles.css"
	htmxPath    = "/static/htmx.1.9.9.min.js"
	prismJSPath = "/static/prism.js"
)

func htmlR(r request.Handler) request.Handler {
	return page.Handler(r, page.Data{
		Title:    "veun-http (demo)",
		CSSPath:  cssPath,
		HTMXPath: htmxPath,
	})
}

var (
	h  = vhttp.Handler
	hf = vhttp.HandlerFunc
)

func server() http.Handler {

	mux := http.NewServeMux()

	// our ajax clicked handler that gets hit by htmx
	mux.Handle("/clicked", h(view.ClickTriggerHandler))

	// components and raw components
	mux.Handle("/component/raw", h(view.ComponentPicker()))
	mux.Handle("/component", h(htmlR(view.ComponentPicker())))

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

	mux.Handle("/lazy", h(htmlR(view.LazyRequestHandler())))

	mux.Handle("/docs", h(htmlR(docsIndex())))
	mux.Handle("/docs/", http.StripPrefix("/docs/", h(htmlR(docsPage()))))

	// mux.Handle("/", handler.OnlyRoot(h(htmlR(view.HomeViewHandler()))))
	mux.Handle("/", handler.OnlyRoot(h(htmlR(index))))

	// Our handler does a 404 fallback between the mux, & static files.
	// we introduce a not-found handler as the last http.Handler to check.
	return handler.Checked(
		mux,                                     // first we we serve our "routes"
		staticFileServer(),                      // if any of those 404 we check static files
		h(htmlR(view.NotFoundRequestHandler())), // if any of those 404 we use our notFound handler
	)
}

var index = request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
	v, err := docsIndexView()
	if err != nil {
		return nil, nil, err
	}

	return veun.Views{
		md.View(md.Bytes(docs.Index)),
		v,
	}, nil, nil
})

func docsIndex() request.Handler {
	return request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
		v, err := docsIndexView()
		if err != nil {
			return nil, nil, err
		}

		return v, nil, nil
	})
}

func docsIndexView() (veun.AsView, error) {
	var filenames []veun.AsView
	if err := fs.WalkDir(docs.DemoServer, ".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".go.md") {
			name := strings.TrimSuffix(path, ".go.md")
			filenames = append(filenames, html.Li(nil, html.A(html.Attrs{
				"href": filepath.Join("/docs", name),
			}, html.Text(name))))
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return html.Ol(nil, filenames...), nil
}

func docsPage() request.Handler {
	return request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {

		if r.URL.Path == "/" {
			return docsIndex().ViewForRequest(r)
		}

		contents, err := docs.DemoServer.ReadFile(r.URL.Path + ".go.md")
		if err != nil {
			slog.Error("failed to read md", "err", err.Error())
			return nil, http.NotFoundHandler(), nil
		}

		return md.View(md.Bytes(contents)), nil, nil
	})
}
