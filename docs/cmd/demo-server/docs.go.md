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
	"github.com/stanistan/veun-http-demo/internal/view/two_column"
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

	var elName veun.AsView
	attrs := el.Attrs{}
	if len(childPages) > 0 {
		name += "/"
		attrs["class"] += " nav-dir"
	} else {
		name += ".md"
	}

	if current == href {
		elName = el.Text(name + " ↞")
		attrs["class"] += " current"
	} else {
		elName = el.A().Attr("href", href).InnerText(name)
	}

	var children []veun.AsView
	for _, name := range childPages {
		children = append(children, el.Li().Content(
			treeView(n.Children[name], current),
		))
	}

	return el.Div().Content(
		el.Div().Attrs(attrs).Content(elName),
		el.Ul().Content(children...),
	)
}
```

## Request Handlers

A fun part here is to use the built-in `request.Handler`, `Always` to make
the actual HTML page we're going to look at for the indexes.

#### Home Page

```go
var index = request.Always(two_column.View{
	Title: "veun-http-demo",
	Nav:   docTree("/"),
	Main:  md.View(static.Index),
})
```

#### Docs Page

Each individual document is mounted at a `/docs/$SLUG` looking
URL, so we can use a request handler specificlaly for that and
mapping it back to the path in our static file server.

So `/docs/$THING` maps to `/docs/$THING.go.md` in our repo, else we 404.

It's would be nice to have our handler also set the title
for the page, and for this we can use `page.DataMutator`.

The view/page takes its own content, lays it out and sets the
title.

```go
func docsForUrl(currentUrl, pathToFile string) (veun.AsView, error) {
	bs, err := static.Docs.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	return two_column.View{
		Nav:   docTree(currentUrl),
		Main:  md.View(bs),
		Title: currentUrl + " | veun-http-demo",
	}, nil
}
```

The handler does the "controller" aspect of making
sure we're rendering the correct view.

Because we are handling both /docs and /docs/ with the same handler,
we strip the prefix here, but keep raw url around because it's useful
for getting the current page.

```go
// FIXME: this specific function/handler should not be looking up the file by path
// but walking the doc tree, and that way if something is a directory, we can
// treat it differently than if if were a file.
var docsHandler = request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
	var (
		rawUrl = r.URL.Path
		url    = strings.TrimPrefix(rawUrl, "/docs")
	)

	page, err := docsForUrl(rawUrl, strings.TrimPrefix(url, "/")+".go.md")
	if err != nil {
		return two_column.View{
			Nav: docTree(rawUrl),
			Main: el.Article().Content(
				el.H2().InnerText(rawUrl),
				el.P().InnerText("pick an .md file"),
                el.Hr(),
			),
			Title: rawUrl + " | veun-http-demo",
		}, nil, nil
	}

	return page, nil, nil
})
```

[md-view]: /docs/internal/view/md/view
