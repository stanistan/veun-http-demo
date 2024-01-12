package view

import (
	"net/http"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"
	"github.com/stanistan/veun/vhttp/request"
)

func NotFound() veun.AsView {
	return el.Div().Content(
		el.P().InnerText("404 Page Not Found."),
		el.P().Content(
			el.A().Attr("href", "/").InnerText("go /"),
		),
	)
}

func NotFoundRequestHandler() request.Handler {
	return request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
		v := NotFound()
		return v, WithStatusCode(http.StatusNotFound), nil
	})
}

func WithStatusCode(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	})
}
