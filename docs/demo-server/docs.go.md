# Rendering our docs

This is one part where the code isn't executed _as is_.
I've embedded the root `/docs/` directory of this repo into
using the `embed` package.

And if we want to render one of these documents we can use
the `md.View()` that we create [here][md-view].

## Requirements

As usual we need to import some things to get this to work.

```go
import (
    "log/slog"
    "path/filepath"
    "net/http"
    "io/fs"
    "strings"
    "sync"

    "github.com/stanistan/veun"
    "github.com/stanistan/veun/html"
    "github.com/stanistan/veun/vhttp/request"

    "github.com/stanistan/veun-http-demo/docs"
    "github.com/stanistan/veun-http-demo/internal/view/md"
)
```

## Docs Index

We want to make a component that lists out all of the docs
that we've written, by their filenames so that we can attach
it to the server.

### Getting the files

In order to do that we need to know what they actually are!
We can use `fs.WalkDir` to figure this out, and generate the entire
list, and do it _once_.

```go
var getDocFiles = sync.OnceValue(
	func() []string {
		var filenames []string
		if err := fs.WalkDir(docs.Docs, ".", func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				slog.Warn("get doc files", "err", err)
				return nil
			}

			if !entry.IsDir() && strings.HasSuffix(path, ".go.md") {
				filenames = append(filenames, path)
			}

			return nil
		}); err != nil {
			panic(err)
		}

		return filenames
	},
)
```

## Index View

Now that we have the files always available, we can make an index page that includes
the directory.

This can start out using the `veun/html` package since we don't _exactly_ know what
we're going to do with this (to render a simple list). But once this gets a bit
more complicated, we can drop it in its own template.


```go
func docFilesIndex() veun.AsView {
	var filenames []veun.AsView

	for _, name := range getDocFiles() {
		var (
			n    = strings.TrimSuffix(name, ".go.md")
			href = filepath.Join("/docs", n)
		)
		filenames = append(
			filenames,
			html.Li(
				nil,
				html.A(html.Attrs{"href": href}, html.Text(n)),
			),
		)
	}

	return html.Section(nil,
        html.H3(nil, html.Text("Articles")),
        html.Ol(nil, filenames...),
    )
}
```

## Request Handlers

A fun part here is to use the built-in `request.Handler`, `Always` to make
the actual HTML page we're going to look at for the indexes.

#### Docs Index

```go
var docsIndex = request.Always(docFilesIndex())
```

#### Home Page

```go
var index = request.Always(veun.Views{
    md.View(docs.Index),
    docFilesIndex(),
})
```

#### Docs Page

Each individual document is mounted at a `/docs/$SLUG` looking
URL, so we can use a request handler specificlaly for that and
mapping it back to the path in our static file server.

So `/docs/$THING` maps to `/docs/$THING.go.md` in our repo, else we 404.

```go
var docsPage = request.HandlerFunc(
	func(r *http.Request) (veun.AsView, http.Handler, error) {
		if r.URL.Path == "" {
			//
			// N.B. if we have no slug do the index!
			return docsIndex.ViewForRequest(r)
		}

		bs, err := docs.Docs.ReadFile(r.URL.Path + ".go.md")
		if err != nil {
			return nil, http.NotFoundHandler(), nil
		}

		return md.View(bs), nil, nil
	},
)
```

[md-view]: /docs/internal/view/md/1-view
