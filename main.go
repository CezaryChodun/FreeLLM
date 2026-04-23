package main

import (
	"github.com/cezarychodun/freellms/app"
	"github.com/cezarychodun/freellms/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
