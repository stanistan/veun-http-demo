# Request Handler

The view that we've built is made for composition,
but so is the request handler!

```go
import (
	"net/http"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/vhttp/request"
)
```

We export a `Handler` function that follows a middleware-like pattern,
taking a handler and returning a handler.

For now given that we do some really simple page construction in the
server this is totally fine, but it can be changed in the future to
take variadic `Option`s do do configuration.

Handler composition is a little verbose (from the function signature),
but really really nice once you're using it in practice.

Once again, you can see that when we're looking at any kind of code
it ends up simply being function/interface composition.

There are small nuances in the way this specific middleware is implemented.

```go
func Handler(rh request.Handler, data Data) request.Handler {
	return request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
```

---

1. It _always_ bubbles up the `next` handler to the caller, regardless of
if there was an error or not.


```go
v, next, err := rh.ViewForRequest(r)
if err != nil {
    return nil, next, err
}
```

2. It fully respects _no view_. If we passed this on to `View`, we would
always be rendering an empty html page and 404s and redirects
would kind of be busted.

```go
if v == nil {
    return nil, next, nil
}
```

3. And if we get something to work with, we wrap with our page, and
see check for that `DataMutator` hook.

```go
return View(v, data), next, nil
```

---

End `Handler`:

```go
	})
}
```

We use this directly in our [`htmlPage` function][html-page].

[html-page]: /docs/demo-server/html-page
