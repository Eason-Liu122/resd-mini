package main

import (
	"embed"
	"resd-mini/core"
)

//go:embed all:web/dist
var assets embed.FS

func main() {
	core.GetApp(assets).Startup()
}
