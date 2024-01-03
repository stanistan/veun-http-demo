package view

import (
	"embed"

	"github.com/stanistan/veun"
)

var (
	//go:embed templates
	templatesFS embed.FS
	templates   = veun.MustParseTemplateFS(templatesFS, "templates/*.tpl")
)
