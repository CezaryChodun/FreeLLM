package main

import (
	"github.com/cezarychodun/freellms/internal/config"
	"github.com/cezarychodun/freellms/internal/modules/app"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
