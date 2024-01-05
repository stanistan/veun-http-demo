package page

import (
	"context"
	_ "embed"

	"github.com/stanistan/veun"
)

func View(v veun.AsView, data Data) view {
	return view{
		body: v,
		data: mutateData(data, v),
	}
}

type view struct {
	body veun.AsView
	data Data
}

func (v view) View(_ context.Context) (*veun.View, error) {
	return veun.V(veun.Template{
		Tpl:   template,
		Slots: veun.Slots{"body": v.body},
		Data:  v.data,
	}), nil
}

var (
	//go:embed template.tpl
	tpl      string
	template = veun.MustParseTemplate("body", tpl)
)
