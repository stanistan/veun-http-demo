You can see an example of the [custom 404 page here](/does_not_exist).


```go
import (
	"net/http"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
	"github.com/stanistan/veun/vhttp/handler"
	"github.com/stanistan/veun/vhttp/request"

	"github.com/stanistan/veun-http-demo/internal/view/doc_tree"
	"github.com/stanistan/veun-http-demo/internal/view/two_column"
)
```

## 404 not found

Repurpose the nav tree and the two-column view to have something that looks ok.

```go
var notFoundView veun.AsView = veun.MustMemo(&two_column.View{
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
var notFoundHandler = request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
	return notFoundView, handler.WithStatus(http.StatusNotFound), nil
})
```
