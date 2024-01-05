package view

import (
	"context"
	"net/http"
	"time"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun-http-demo/internal/view/page"
	"github.com/stanistan/veun/vhttp/request"
)

func HomeViewHandler() request.Handler {
	return request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
		if r.URL.Query().Get("sleep") != "" {
			time.Sleep(2 * time.Second)
		}
		return home{
			Eager: r.URL.Query().Get("fast") != "",
		}, nil, nil
	})
}

type home struct {
	Eager bool
}

func (v home) components() veun.AsView {
	if v.Eager {
		return DefinedComponents
	} else {
		return Lazy{Endpoint: "/components", Delay: "2s"}
	}
}

func (v home) View(_ context.Context) (*veun.View, error) {
	return veun.V(veun.Template{
		Tpl:   templates.Lookup("home.tpl"),
		Slots: veun.Slots{"components": v.components()},
	}), nil
}

func (m home) SetPageData(d *page.Data) {
	d.Title = "Home"
}
