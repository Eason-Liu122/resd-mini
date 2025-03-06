package main

import (
	"embed"
	"log"
	"resd-mini/core"
)

//go:embed all:web/dist
var assets embed.FS

func main() {
	logo := `
  ____  _____ ____  ____            __  __ ___ _   _ ___ 
 |  _ \| ____/ ___||  _ \          |  \/  |_ _| \ | |_ _|
 | |_) |  _| \___ \| | | |  _____  | |\/| || ||  \| || | 
 |  _ <| |___ ___) | |_| | |_____| | |  | || || |\  || | 
 |_| \_\_____|____/|____/          |_|  |_|___|_| \_|___|`

	log.Println(logo)
	core.GetApp(assets).Startup()
}
