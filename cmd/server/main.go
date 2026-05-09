package main

import (
	"github.com/cezarychodun/freellms/internal/app"
	"github.com/cezarychodun/freellms/internal/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
