package view

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
	"github.com/stanistan/veun/vhttp/request"
)

func LazyRequestHandler() request.Handler {
	return request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
		return Lazy{
			Endpoint: "/home?sleep=true",
		}, nil, nil
	})
}

type Lazy struct {
	Endpoint    string
	Placeholder veun.AsView
	Delay       string
	UseTpl      bool
}

func (v Lazy) View(ctx context.Context) (*veun.View, error) {
	if v.Endpoint == "" {
		return nil, fmt.Errorf("needs endpoint")
	}

	var placeholder veun.AsView
	if v.Placeholder == nil {
		placeholder = el.Em().InnerText("...loading...")
	} else {
		placeholder = v.Placeholder
	}

	if !v.UseTpl {
		return el.Div().
			Attrs(el.Attrs{
				"hx-get":     v.Endpoint,
				"hx-trigger": "load delay:" + v.Delay,
			}).
			Content(placeholder).
			View(ctx)
	}

	return veun.V(veun.Template{
		Tpl:   templates.Lookup("lazy.tpl"),
		Data:  v,
		Slots: veun.Slots{"placeholder": placeholder},
	}), nil
}

func (v Lazy) Title() string {
	var title string
	if v.UseTpl {
		title = "using lazy.tpl"
	} else {
		title = "using Div"
	}

	return fmt.Sprintf("%s (delay=%s)", title, v.Delay)
}
