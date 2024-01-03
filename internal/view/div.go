package view

import (
	"context"
	"strings"
	"text/template"

	"github.com/stanistan/veun"
)

type Attrs map[string]string

func (a Attrs) InTag(name string) string {
	var sb strings.Builder

	sb.WriteString("<")
	sb.WriteString(name)
	sb.WriteString(" ")

	for k, v := range a {
		template.HTMLEscape(&sb, []byte(k))
		sb.WriteString(`="`)
		template.HTMLEscape(&sb, []byte(v))
		sb.WriteString(`" `)
	}

	sb.WriteString(">")

	return sb.String()
}

type Props struct {
	Attrs   Attrs
	Content veun.AsView
}

type Element struct {
	Name string
	Props
}

func (e Element) View(ctx context.Context) (*veun.View, error) {
	return veun.V(veun.Views{
		veun.Raw(e.Attrs.InTag(e.Name)),                 // opening tag
		e.Content,                                       // content
		veun.Raw("</"), veun.Raw(e.Name), veun.Raw(">"), // closing tag
	}), nil
}

func Div(r veun.AsView, attrs Attrs) veun.AsView {
	return Element{
		Name:  "div",
		Props: Props{Attrs: attrs, Content: r},
	}
}

func Li(r veun.AsView, attrs Attrs) veun.AsView {
	return Element{
		Name:  "li",
		Props: Props{Attrs: attrs, Content: r},
	}
}

func Ul(attrs Attrs, r ...veun.AsView) veun.AsView {
	return Element{
		Name:  "li",
		Props: Props{Attrs: attrs, Content: veun.Views(r)},
	}
}

func el(name string, attrs Attrs, children ...veun.AsView) veun.AsView {
	return Element{
		Name: name,
		Props: Props{
			Attrs:   attrs,
			Content: veun.Views(children),
		},
	}
}
