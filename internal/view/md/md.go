package md

import (
	"bytes"
	"context"
	"embed"
	"html/template"

	"github.com/stanistan/veun"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Source interface {
	Source() ([]byte, error)
}

type SourceFunc func() ([]byte, error)

func (f SourceFunc) Source() ([]byte, error) {
	return f()
}

type Bytes []byte

func (b Bytes) Source() ([]byte, error) {
	return b, nil
}

func ReadFile(f embed.FS, name string) Source {
	return SourceFunc(func() ([]byte, error) {
		return f.ReadFile(name)
	})
}

func View(s Source) view {
	return view{
		source: s,
	}
}

type view struct {
	source Source
}

func (v view) AsHTML(_ context.Context) (template.HTML, error) {
	in, err := v.source.Source()
	if err != nil {
		return template.HTML(""), err
	}

	var out bytes.Buffer
	err = md.Convert(in, &out)
	if err != nil {
		return template.HTML(""), err
	}

	return template.HTML(out.String()), nil
}

func (v view) View(_ context.Context) (*veun.View, error) {
	return veun.V(v), nil
}

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)
