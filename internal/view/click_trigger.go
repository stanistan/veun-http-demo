package view

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/vhttp"
)

var ClickTriggerHandler = vhttp.HandlerFunc(
	func(r *http.Request) (veun.AsView, http.Handler, error) {

		count, err := strconv.Atoi(r.URL.Query().Get("count"))
		if err != nil {
			return nil, nil, errors.New("count is not a number")
		}

		return ClickTrigger{Count: count}, nil, nil
	},
)

type ClickTrigger struct {
	Count int
}

func (v ClickTrigger) View(_ context.Context) (*veun.View, error) {
	return veun.V(veun.Template{
		Tpl:  templates.Lookup("click_trigger.tpl"),
		Data: v.Count + 1,
	}), nil
}

func (v ClickTrigger) Title() string {
	return "HTMX enabled click counter"
}
