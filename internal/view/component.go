package view

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
)

type Component interface {
	veun.AsView
	Title() string
}

type Components []Component

func ComponentView(c Component) component {
	return component{
		Name: fmt.Sprintf("%T - %s", c, c.Title()),
		Body: c,
	}
}

func (v Components) View(ctx context.Context) (*veun.View, error) {
	items := make(veun.Views, len(v))
	for idx, c := range v {
		items[idx] = veun.Views{
			ComponentView(c),
			ComponentLink(idx),
		}
	}

	return veun.V(items), nil
}

type component struct {
	Name      string
	Body      veun.AsView
	BodyClass string
}

func (v component) template() veun.Template {
	return veun.Template{
		Tpl:   templates.Lookup("component.tpl"),
		Slots: veun.Slots{"body": v.Body},
		Data:  v,
	}
}

func (v component) View(_ context.Context) (*veun.View, error) {
	// N.B. every single component gets error handling.
	// so no single one gets to break anything else.
	return veun.V(v.template()).WithErrorHandler(v), nil
}

func (v component) ViewForError(ctx context.Context, err error) (veun.AsView, error) {
	return component{
		Name: v.Name,
		Body: veun.Views{
			el.Div().Content(el.Strong().InnerText("Error Captured:")),
			niceError(err),
		},
		BodyClass: "error",
	}, nil
}

func niceError(err error) veun.AsView {
	var (
		errs = []string{err.Error()}
	)

	for {
		err = errors.Unwrap(err)
		if err == nil {
			break
		}

		var (
			msg     = err.Error()
			prevIdx = len(errs) - 1
			prev    = errs[prevIdx]
		)

		errs[prevIdx] = strings.TrimRight(
			prev[:len(prev)-len(msg)], ": ",
		)

		errs = append(errs, msg)
	}

	lis := make(veun.Views, len(errs))
	for idx, err := range errs {
		lis[idx] = el.Li().InnerText(err)
	}

	return el.Ul().Content(lis)
}
