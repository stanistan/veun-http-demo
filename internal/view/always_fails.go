package view

import (
	"context"
	"errors"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/html"
)

type AlwaysFails struct {
	OwnErrorCapture bool
}

func (v AlwaysFails) Title() string {
	if v.OwnErrorCapture {
		return "captures itself"
	}

	return "captured by component"
}

func (v AlwaysFails) View(_ context.Context) (*veun.View, error) {
	if !v.OwnErrorCapture {
		return nil, errors.New("this view will always fail")
	}

	// N.B. Why yes, this is a recursive definition.
	return veun.V(AlwaysFails{}).WithErrorHandler(v), nil
}

func (v AlwaysFails) ViewForError(ctx context.Context, err error) (veun.AsView, error) {
	return veun.Views{
		html.Strong(nil, html.Text("Error Inline:")),
		niceError(err),
	}, nil
}
