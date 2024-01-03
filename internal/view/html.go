package view

import (
	"context"
	"net/http"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/vhttp/request"
)

func HTML(rh request.Handler, data HTMLData) request.Handler {
	return request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
		v, next, err := rh.ViewForRequest(r)
		if err != nil {
			return nil, next, err
		}

		if v == nil {
			return nil, next, nil
		}

		return htmlView(v, data), next, nil
	})
}

func htmlView(v veun.AsView, data HTMLData) html {
	if mutator, ok := v.(HTMLDataMutator); ok {
		mutator.SetHTMLData(&data)
	}

	return html{
		Body:     v,
		HTMLData: data,
	}
}

type HTMLDataMutator interface {
	SetHTMLData(data *HTMLData)
}

type HTMLData struct {
	Title    string
	CSSPath  string
	HTMXPath string
}

type html struct {
	Body veun.AsView

	HTMLData
}

func (v html) View(_ context.Context) (*veun.View, error) {
	return veun.V(veun.Template{
		Tpl:   templates.Lookup("html.tpl"),
		Data:  v.HTMLData,
		Slots: veun.Slots{"body": v.Body},
	}), nil
}
