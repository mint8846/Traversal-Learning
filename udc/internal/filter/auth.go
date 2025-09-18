package filter

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mint8846/Traversal-Learning/udc/internal/config"
)

type AuthMiddleware struct {
	cfg *config.Config
}

func New(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		cfg: cfg,
	}
}

func (m *AuthMiddleware) Verify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		clientIP := c.RealIP()

		if m.cfg.ODC.IP != "" && m.cfg.ODC.IP != clientIP {
			log.Printf("Access denied: IP %s is not allowed (expected: %s)", clientIP, m.cfg.ODC.IP)
			return echo.NewHTTPError(http.StatusForbidden, "IP not allowed")
		}

		containerId := c.Request().Header.Get("X-Container-ID")
		if m.cfg.ODC.ID != "" && m.cfg.ODC.ID != containerId {
			log.Printf("Access denied: ID %s is not allowed (expected: %s)", containerId, m.cfg.ODC.ID)
			return echo.NewHTTPError(http.StatusForbidden, "ID not allowed")
		}

		return next(c)
	}
}
