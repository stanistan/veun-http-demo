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

	"github.com/stanistan/veun-http-demo/docs"
	"github.com/stanistan/veun-http-demo/internal/components"
	"github.com/stanistan/veun-http-demo/internal/view/doc_tree"
	"github.com/stanistan/veun-http-demo/internal/view/md"
	"github.com/stanistan/veun-http-demo/internal/view/title"
	"github.com/stanistan/veun-http-demo/internal/view/two_column"
)
```

## Request Handlers

A fun part here is to use the built-in `request.Handler`, `Always` to make
the actual HTML page we're going to look at for the indexes.

#### Home Page

```go
var index = request.Always(&two_column.View{
    Title: "veun-http-demo",
    Nav:   doc_tree.View("/"),
    Main:  el.Article().Content(md.View(docs.Index)),
})
```

#### Docs Page

Each individual document is mounted at a `/docs/$SLUG` looking
URL, so we can use a request handler specificlaly for that and
mapping it back to the path in our static file server.

So `/docs/$THING.md` maps to `/docs/$THING.go.md` in our repo.

Our handler can also set the title for the page using `page.DataMutator`.

```go
// FIXME: this specific function/handler should not be looking up the file by path
// but walking the doc tree, and that way if something is a directory, we can
// treat it differently than if if were a file.
var docsHandler = request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
	var (
		rawUrl = r.URL.Path
		url    = strings.TrimPrefix(rawUrl, "/docs")
		docUrl = strings.TrimSuffix(strings.TrimPrefix(url, "/"), ".md") + ".go.md"

		next http.Handler
	)

	content := docPageContent(rawUrl, docUrl)
	if content == nil {
        return nil, http.NotFoundHandler(), nil
	}

	return &two_column.View{
		Nav:   doc_tree.View(rawUrl),
		Main:  content,
		Title: rawUrl + " | veun-http-demo",
	}, next, nil
})
```

The page content itself reads the file from the static documentation
and attempts to render it into markdown.

Importantly, if the file also has an associated [component](/docs/internal/components/registry.md),
then that gets prepended to the content as well.

```go
func docPageContent(currentUrl, pathToFile string) veun.AsView {
	bs, err := docs.Docs.ReadFile(pathToFile)
	if err != nil {
        return nil
	}

    content := md.View(bs)
    if component, ok := components.ForFullURL(currentUrl); ok {
        content = veun.Views{component, content}
    }

    return el.Article().Content(
        title.View(currentUrl),
        el.Hr(),
        content,
    )
}
```


[md-view]: /docs/internal/view/md/view.md
