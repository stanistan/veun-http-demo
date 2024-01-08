package docs

import (
	"embed"
)

var (
	//go:embed demo-server internal
	Docs embed.FS

	//go:embed index.md
	Index []byte
)
