package docs

import (
	"embed"
)

var (
	//go:embed cmd internal
	Docs embed.FS

	//go:embed index.md
	Index []byte
)
