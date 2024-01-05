package page

import (
	"net/http"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/vhttp/request"
)

func Handler(rh request.Handler, data Data) request.Handler {
	return request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
		v, next, err := rh.ViewForRequest(r)
		if err != nil {
			return nil, next, err
		}

		if v == nil {
			return nil, next, nil
		}

		return View(v, data), next, nil
	})
}
