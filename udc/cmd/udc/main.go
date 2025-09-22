package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mint8846/Traversal-Learning/udc/internal/config"
	"github.com/mint8846/Traversal-Learning/udc/internal/router"
	"github.com/mint8846/Traversal-Learning/udc/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		return
	}

	e := echo.New()

	service.Initialize(cfg)
	router.Setup(e, cfg)

	e.Logger.Fatal(e.Start(cfg.Port))
}
