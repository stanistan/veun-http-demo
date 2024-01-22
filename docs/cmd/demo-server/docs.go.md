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
	"github.com/stanistan/veun-http-demo/internal/components"
	"github.com/stanistan/veun-http-demo/internal/docs"
	"github.com/stanistan/veun-http-demo/internal/view/md"
	"github.com/stanistan/veun-http-demo/internal/view/two_column"
)
```

## Docs Index

We want to make a component that lists out all of the docs
that we've written, by their filenames so that we can attach
it to the server.

### Getting the tree

This is in the [`internal/docs`](/docs/internal/docs/tree.md) package.

## Index View

Now that we have the files always available, we can make an index page that includes
the directory.

We pass in the current url so we can set an active node on the tree.

```go
func docTree(current string) veun.AsView {
	return el.Div().Class("doc-tree").Content(
		treeView(docs.Tree(), current),
	)
}
```

And our tree view function. This is recursive and walks the
entire tree to build out the nav.

```go
func treeView(n docs.Node, current string) veun.AsView {
	childPages := n.SortedKeys()
	name, href := n.LinkInfo()

	var elName veun.AsView
	attrs := el.Attrs{}
	if len(childPages) > 0 {
		attrs["class"] += " nav-dir"
	}

	if current == href {
		elName = el.Text(name + " ↞")
		attrs["class"] += " current"
	} else {
		elName = el.A().Attr("href", href).InnerText(name)
	}

	var childContent veun.AsView
	if len(childPages) > 0 {
		var children []veun.AsView
		for _, name := range childPages {
			children = append(children, el.Li().Content(
				treeView(n.Children[name], current),
			))
		}
		childContent = el.Ul().Content(children...)
	}

	return el.Div().Content(
		el.Div().Attrs(attrs).Content(elName),
		childContent,
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
	Main:  el.Article().Content(md.View(static.Index)),
})
```

#### Docs Page

Each individual document is mounted at a `/docs/$SLUG` looking
URL, so we can use a request handler specificlaly for that and
mapping it back to the path in our static file server.

So `/docs/$THING` maps to `/docs/$THING.go.md` in our repo.

Our handler can also set the title for the page using `page.DataMutator`.

```go
func docPageContent(currentUrl, pathToFile string) veun.AsView {
	bs, err := static.Docs.ReadFile(pathToFile)
	if err != nil {
        return fallbackContent(currentUrl)
	}

    content := md.View(bs)

    if component, ok := components.ForURL(
        strings.TrimPrefix(currentUrl, "/docs/internal/components/"),
    ); ok {
        content = veun.Views{component, content}
    }

    return el.Article().Content(
        el.H1().InnerText(currentUrl),
        el.Hr(),
        content,
    )
}

func fallbackContent(url string) veun.AsView {
	return el.Article().Content(
		el.H1().InnerText(url),
		el.P().InnerText("this is fallback content."),
		el.Hr(),
		el.P().InnerText("probably for a directory"),
	)
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

    content := docPageContent(
        rawUrl,
        strings.TrimSuffix(strings.TrimPrefix(url, "/"), ".md")+".go.md",
    )

	return two_column.View{
		Nav:   docTree(rawUrl),
		Main:  content,
		Title: rawUrl + " | veun-http-demo",
	}, nil, nil
})
```

[md-view]: /docs/internal/view/md/view.md
