// Code Generated by github.com/stanistan/veun-http-demo/cmd/lit-gen; DO NOT EDIT.
package main

import (
	"net/http"

	"github.com/stanistan/veun"
	"github.com/stanistan/veun/html"
	"github.com/stanistan/veun/vhttp/request"

	static "github.com/stanistan/veun-http-demo/docs"
	"github.com/stanistan/veun-http-demo/internal/docs"
	"github.com/stanistan/veun-http-demo/internal/view/md"
)

func treeView(n docs.Node) veun.AsView {
	var children []veun.AsView
	for _, name := range n.Sorted() {
		children = append(children, html.Li(nil, treeView(n.Children[name])))
	}

	name, href := n.LinkInfo()
	if len(children) == 0 {
		return html.Div(nil, html.A(html.Attrs{"href": href}, html.Text(name)))
	}

	return html.Div(nil,
		html.Div(nil, html.Text(name+"/")),
		html.Ul(nil, children...))
}

func docFilesIndex() veun.AsView {
	return html.Div(html.Attrs{"class": "doc-tree"}, treeView(docs.Tree()))
}

var docsIndex = request.Always(docFilesIndex())

var index = request.Always(veun.Views{
	md.View(static.Index),
	veun.Raw("<hr />"),
	docFilesIndex(),
})

var docsPage = request.HandlerFunc(func(r *http.Request) (veun.AsView, http.Handler, error) {
	if r.URL.Path == "" {
		return docsIndex.ViewForRequest(r)
	}

	bs, err := static.Docs.ReadFile(r.URL.Path + ".go.md")
	if err != nil {
		return nil, http.NotFoundHandler(), nil
	}

	return html.Div(
		html.Attrs{"class": "doc-page-cols"},
		html.Div(nil, docFilesIndex()),
		html.Div(nil, md.View(bs)),
	), nil, nil
})
