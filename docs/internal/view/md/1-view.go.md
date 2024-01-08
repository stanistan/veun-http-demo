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
	"github.com/stanistan/veun/html"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)
```

## Goldmark Configuration

It's actually _very very simple_ to make a trivial wrapper
around goldmark.

Similar to `veun.Template` we want something that takes our
input, and `goldmark` uses `bytes`, and the configuration
of the library and does what it needs to.

For our use case, we can make a singleton, and if we really
want to at some point in the future, make this configurable.

```go
var markdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)
```

## Our HTMLRenderable

This _does_ take a `goldmark.Markdown` as an option, it's a
passthrough struct of input and md converter to produce HTML,
that's all it takes to make a `veun.HTMLRenderable`.

```go
type Markdown struct {
	Bytes    []byte
	Markdown goldmark.Markdown
}

func (m Markdown) AsHTML(_ context.Context) (template.HTML, error) {
	var out bytes.Buffer
	if err := m.Markdown.Convert(m.Bytes, &out); err != nil {
		return template.HTML(""), err
	}

	return template.HTML(out.String()), nil
}
```

## Our view

Our view uses just the bytes, and the singleton `markdown`.

```go
type view struct {
	Bytes []byte
}

func (v view) View(_ context.Context) (*veun.View, error) {
	return veun.V(Markdown{
		Bytes:    v.Bytes,
		Markdown: markdown,
	}), nil
}
```

There was another option to have the `view` take _sources_
that can fail, but and produce bytes as output, but this
is fine, and dumb, and simple, and simple is good.

## Constructor

This is trivial, and allows us to have more semantic
HTML around our generated markdown-- like an `<article>`.

```go
func View(bs []byte) veun.AsView {
	return html.Article(nil, view{Bytes: bs})
}
```


[goldmark]: https://github.com/yuin/goldmark
