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

	authMiddleware := filter.New(cfg)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("auth_middleware", authMiddleware)
			return next(c)
		}
	})

	// API route group
	api := e.Group("/api")

	api.GET("/connect", handler.Connect, authMiddleware.CreateSession)
	api.POST("/model", handler.Model, authMiddleware.Authenticate)
	api.POST("/result", handler.Result, authMiddleware.AuthenticateWithCleanup)
}
