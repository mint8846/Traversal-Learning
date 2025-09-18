package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mint8846/Traversal-Learning/udc/internal/config"
	"github.com/mint8846/Traversal-Learning/udc/internal/filter"
	"github.com/mint8846/Traversal-Learning/udc/internal/handler"
)

func Setup(e *echo.Echo, cfg *config.Config) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// defer func()
	handler.SetId(cfg.ID)

	authMiddleware := filter.New(cfg)

	// API route group
	api := e.Group("/api")
	api.Use(authMiddleware.Verify)

	api.GET("/connect", handler.Connect)
	api.POST("/model", handler.Model)
	api.POST("/result", handler.Result)
}
