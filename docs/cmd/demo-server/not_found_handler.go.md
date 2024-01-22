```go
import (
	"context"
	"net/http"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
	"github.com/stanistan/veun/vhttp/request"

	"github.com/stanistan/veun-http-demo/internal/view/doc_tree"
	"github.com/stanistan/veun-http-demo/internal/view/two_column"
)
```

## 404 not found

Repurpose the nav tree and the two-column view to have something that looks ok.

```go
var notFoundView veun.AsView = mustMemo(&two_column.View{
	Title: "404 Not Found",
	Nav:   doc_tree.View(""),
	Main: el.Article().Content(
		el.H1().InnerText("404 Not Found"),
		el.Hr(),
		el.A().Attr("href", "/").InnerText("/ go home").In(el.P()),
		el.P().InnerText("The content you were looking for was not found."),
	),
})
```

This view is memoized and rendered when the server is started, but continues
to be represented as a view in the typesystem.

---

The handler returns both the view and a handler that _only_ sets the
response code to 404.

```go
func notFoundHandler() http.Handler {
	return h(html(request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
        return notFoundView, withStatus(http.StatusNotFound), nil
	})))
}
```

And we have a helper middleware here that writes the status code.

```go
func withStatus(code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	})
}
```

### memoization

In order to make a cached view (in memory), we need to render one and wrap it
as a `Raw` type.

```go
func mustMemo(v veun.AsView) veun.Raw {
    out, err := veun.Render(context.Background(), v)
    if err != nil {
        panic(err)
    }

    return veun.Raw(out)
}
```
