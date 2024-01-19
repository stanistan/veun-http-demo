package view

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
	"github.com/stanistan/veun/vhttp/request"
)

func ComponentPicker() request.Handler {
	return request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
		// get our idx
		idx, err := strconv.Atoi(r.URL.Query().Get("idx"))
		if err != nil {
			return nil, http.NotFoundHandler(), nil
		}

		if idx > len(DefinedComponentHandlers)-1 {
			return nil, http.NotFoundHandler(), nil
		}

		h := DefinedComponentHandlers[idx]
		return h.ViewForRequest(r)
	})
}

var DefinedComponents = Components{
	ClickTrigger{},
	Delazy("2s", false),
	Delazy("8s", true),
}

var DefinedComponentHandlers = []request.Handler{
	ComponentHandler(ClickTriggerHandler),
	DelazyHandler("2s", false),
	DelazyHandler("8s", true),
}

func Delazy(delay string, tpl bool) Component {
	return Lazy{
		Endpoint: "/v/echo?in=" + url.QueryEscape(fmt.Sprintf("After a %s delay", delay)),
		Delay:    delay,
		UseTpl:   tpl,
		Placeholder: veun.Views{
			el.Em().InnerText("...loading..."),
			el.Text(fmt.Sprintf(" incurring a %s delay", delay)),
		},
	}
}

func DelazyHandler(delay string, tpl bool) request.Handler {
	return ComponentHandler(request.Always(Delazy(delay, tpl)))
}

func ComponentHandler(h request.Handler) request.Handler {
	return request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
		v, next, err := h.ViewForRequest(r)
		if v == nil || err != nil {
			return nil, next, err
		}

		c, ok := v.(Component)
		if !ok {
			return nil, nil, fmt.Errorf("expected a component got %T", v)
		}

		return ComponentView(c), next, nil
	})
}
