package view

import (
	"net/http"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/html"
	"github.com/stanistan/veun/vhttp/request"
)

func NotFound() veun.AsView {
	return html.Div(
		nil,
		html.P(nil, html.Text("404 Page Not Found.")),
		html.P(nil, html.A(html.Attrs{"href": "/"}, html.Text("go /"))),
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
