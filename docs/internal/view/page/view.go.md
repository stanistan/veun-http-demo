# page

This defines the page structure: html element, title, css, js files, etc.

This is one we can use both `veun/template` and `veun/el`.

For this, we opt to use go templating directly.

```go
import (
	"context"
	_ "embed"

	"github.com/stanistan/veun"
	t "github.com/stanistan/veun/template"
)

//go:embed template.tpl
var tpl string
var template = t.MustParse("page", tpl)
```

You can see the actual template [here][template-link].

## Data

We see in the `main` function, that declare the assets (css/js) for the
view _external_ to it, it's pretty non-functional on its own.

```go
type Data struct {
	Title    string
	CSSFiles []string
	JSFiles  []string
	IsMobile bool
}
```

## The View

Now that we have the html defined, let's give it an actual view. This one
is really simple, just passing through the data and the view for the body
slot.

```go
type view struct {
	body veun.AsView
	data Data
}

func (v view) View(_ context.Context) (*veun.View, error) {
	return veun.V(t.Template{
		Tpl:   template,
		Slots: t.Slots{"body": v.body},
		Data:  v.data,
	}), nil
}
```

The view struct itself is very small, and we can add more things
to it as we go, but let's take a closer tlook at that _data_.

### Mutation hooks

A thing we provide specifically for the `page.Data` is an option for
the view that gets wrapped by the page to provide some configuration
back up the tree (optionally).

We do this by giving the view an interface to fullfill.

```go
type DataMutator interface {
	SetPageData(d *Data)
}
```

And the function we use to for construction will invoke this for us.

```go
func View(v veun.AsView, data Data) view {
	return view{
		body: v,
		data: mutateData(data, v),
	}
}

func mutateData(d Data, with any) Data {
	if m, ok := with.(DataMutator); ok {
		m.SetPageData(&d)
	}

	return d
}
```

[template-link]: https://github.com/stanistan/veun-http-demo/blob/main/internal/view/page/template.tpl
