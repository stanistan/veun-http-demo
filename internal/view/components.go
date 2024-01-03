package view

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/vhttp/request"
)

var PickAComponent = request.HandlerFunc(
	func(r *http.Request) (veun.AsView, http.Handler, error) {
		idx, err := strconv.Atoi(r.URL.Query().Get("idx"))
		if err != nil {
			return nil, http.NotFoundHandler(), nil
		}

		if idx > len(DefinedComponents)-1 {
			return nil, http.NotFoundHandler(), nil
		}

		return ComponentView(DefinedComponents[idx]), nil, nil
	},
)

func ComponentLink(idx int) veun.AsView {
	return el("div", Attrs{"class": "component-permalink"},
		el("a", Attrs{"href": fmt.Sprintf("/component?idx=%d", idx)},
			veun.Raw("premalink"),
		),
	)
}

var DefinedComponents = Components{
	ClickTrigger{},
	Delazy("2s", false),
	Delazy("8s", true),
	AlwaysFails{},
	AlwaysFails{true},
}

func Delazy(delay string, tpl bool) Component {
	return Lazy{
		Endpoint: "/v/echo?in=" + url.QueryEscape(fmt.Sprintf("After a %s delay", delay)),
		Delay:    delay,
		UseTpl:   tpl,
		Placeholder: veun.Views{
			veun.Raw("<em>...loading...</em>"),
			veun.Raw(" "),
			veun.Raw(fmt.Sprintf("incurring a %s delay", delay)),
		},
	}
}
