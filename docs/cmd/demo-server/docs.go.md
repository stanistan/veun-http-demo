# Rendering our docs

This is one part where the code isn't executed _as is_.
I've embedded the root `/docs/` directory of this repo into
using the `embed` package.

And if we want to render one of these documents we can use
the `md.View()` that we create [here][md-view].

```go
import (
    "net/http"
    "strings"

    "github.com/stanistan/veun"
    "github.com/stanistan/veun/el"
    "github.com/stanistan/veun/vhttp/request"

    static "github.com/stanistan/veun-http-demo/docs"
    "github.com/stanistan/veun-http-demo/internal/docs"
    "github.com/stanistan/veun-http-demo/internal/view/md"
)
```

## Docs Index

We want to make a component that lists out all of the docs
that we've written, by their filenames so that we can attach
it to the server.

### Getting the files

This is in the [`internal/docs`](/docs/internal/docs/tree) package.

## Index View

Now that we have the files always available, we can make an index page that includes
the directory.

This can start out using the `veun/html` package since we don't _exactly_ know what
we're going to do with this (to render a simple list). But once this gets a bit
more complicated, we can drop it in its own template.

We pass in the current url so we can set an active node on the tree.

```go
func docTree(current string) veun.AsView {
	return el.Div().Class("doc-tree").Content(
		treeView(docs.Tree(), current),
	)
}
```

And our tree view, this could be a struct instead of a function.

```go
func treeView(n docs.Node, current string) veun.AsView {
	childPages := n.SortedKeys()
	name, href := n.LinkInfo()

	if len(childPages) == 0 {
		if current == href {
			return el.Div().Class("current").InnerText(name)
		} else {
			return el.Div().Content(
                el.A().Attr("href", href).InnerText(name),
            )
		}
	}

	var children []veun.AsView
	for _, name := range childPages {
		children = append(children, el.Li().Content(
			treeView(n.Children[name], current),
		))
	}

	return el.Div().Content(
		el.Div().InnerText(name+"/"),
		el.Ul().Content(children...),
	)
}
```

## Request Handlers

A fun part here is to use the built-in `request.Handler`, `Always` to make
the actual HTML page we're going to look at for the indexes.

#### Home Page

```go
var index = request.Always(veun.Views{
    md.View(static.Index),
    veun.Raw("<hr />"),
    docTree("/"),
})
```

#### Docs Page

Each individual document is mounted at a `/docs/$SLUG` looking
URL, so we can use a request handler specificlaly for that and
mapping it back to the path in our static file server.

So `/docs/$THING` maps to `/docs/$THING.go.md` in our repo, else we 404.

```go
var docsHandler = request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
    url := strings.TrimPrefix(r.URL.Path, "/docs")
	if url == "/" || url == "" {
        return docTree("/"), nil, nil
	}

    url = strings.TrimPrefix(url, "/")
	bs, err := static.Docs.ReadFile(url+ ".go.md")
	if err != nil {
		return nil, http.NotFoundHandler(), nil
	}

    return el.Div().Class("doc-page-cols").Content(
        el.Div().Content(docTree(r.URL.Path)),
        el.Div().Content(md.View(bs)),
    ), nil, nil
})
```

[md-view]: /docs/internal/view/md/view
