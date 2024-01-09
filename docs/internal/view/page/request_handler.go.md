# Request Handler

The view that we've built is made for composition,
but so is the request handler!

We use this directly in our [`htmlPage` function][html-page].

## Our dependencies

Pretty standard for something dealing with an http request.

```go
import (
	"net/http"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/vhttp/request"
)
```

## The Handler

We export a `Handler` function that follows a middleware-like pattern,
taking a handler and returning a handler. Handler composition is a little
verbose (from the function signature), but really really nice once you're
using it in practice.

```go
func Handler(rh request.Handler, data Data) request.Handler {
	return request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
```

---

Once again, you can see that when we're looking at any kind of code
it ends up simply being function/interface composition.

There are small nuances in the way this specific middleware is implemented.

### Errors and handlers

It _always_ bubbles up the `next` handler to the caller, regardless of
if there was an error or not.


```go
v, next, err := rh.ViewForRequest(r)
if err != nil {
    return nil, next, err
}
```

### Middlewares

It fully respects _no view_. If we passed this on to `View`, we would
always be rendering an empty html page and 404s and redirects
would kind of be busted. This also allows for a `request.Handler` to
produce _anything_.

```go
if v == nil {
    return nil, next, nil
}
```

### Nice defaults

And if we get something to work with, we wrap with our page, and
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


[html-page]: /docs/demo-server/html-page
