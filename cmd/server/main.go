package main

import (
	"fmt"

	"github.com/cezarychodun/freellms/internal/app"
	"github.com/cezarychodun/freellms/internal/config"
)

func main() {
	cfg := config.GetConfig()

	a := &app.App{}
	a.Initialize(cfg)
	a.Run(fmt.Sprintf(":%d", cfg.Port))
}
