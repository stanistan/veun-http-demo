# Rendering Markdown

Golang has a really good markdown rendering library called
[goldmark][goldmark], and we're going to use it in conjunction
with `veun.View` to make our pages.

```go
import (
	"bytes"
	"context"
	"html/template"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)
```

## Goldmark Configuration

It's actually _very very simple_ to make a trivial wrapper
around goldmark.

Similar to `veun/template.Template` we want something that takes our
input, and `goldmark` uses `bytes`, and the configuration
of the library and does what it needs to.

For our use case, we can make a singleton, and if we really
want to at some point in the future, make this configurable.

```go
var md = goldmark.New(goldmark.WithExtensions(extension.GFM))
```

## Our HTMLRenderable

This takes a `goldmark.Markdown` as an struct parameter, and
defaults to our global `md` if it wasn't set.

That's all it takes to make a `veun.HTMLRenderable`.

```go
type view struct {
	bytes []byte
	md    goldmark.Markdown
}

func (v view) AsHTML(_ context.Context) (template.HTML, error) {
	var r goldmark.Markdown
	if v.md != nil {
		r = v.md
	} else {
		r = md
	}

	var out bytes.Buffer
	if err := r.Convert(v.bytes, &out); err != nil {
		var empty template.HTML
		return empty, err
	}

	return template.HTML(out.String()), nil
}
```

## Our view

Our view uses just the bytes, and the singleton `markdown`.

```go
func (v view) View(_ context.Context) (*veun.View, error) {
    return veun.V(v), nil
}
```

There was another option to have the `view` take _sources_
that can fail, but and produce bytes as output, but this
is fine, and dumb, and simple, and simple is good.

## Constructor

This is trivial, and allows us to have more semantic
HTML around our generated markdown.

```go
func View(bs []byte) veun.AsView {
	return el.Div().Attr("class", "md").Content(
		view{bytes: bs},
	)
}
```

[goldmark]: https://github.com/yuin/goldmark
