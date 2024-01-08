package docs

import (
	"embed"
)

var (
	//go:embed demo-server
	DemoServer embed.FS

	//go:embed internal
	Internal embed.FS

	//go:embed index.md
	Index []byte
)
