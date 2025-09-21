package filter

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mint8846/Traversal-Learning/udc/internal/config"
	"github.com/mint8846/Traversal-Learning/udc/internal/utils"
)

type Session struct {
	ID        string
	CreatedAt time.Time
}

type AuthMiddleware struct {
	cfg      *config.Config
	sessions map[string]*Session
	mutex    sync.RWMutex
}

const (
	SessionKey = "session_key"
)

func New(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:      cfg,
		sessions: make(map[string]*Session),
		mutex:    sync.RWMutex{},
	}
}

func (a *AuthMiddleware) CreateSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := a.validateClient(c); err != nil {
			return err
		}

		err := next(c)
		if err != nil {
			return err
		}

		if err = a.createSession(c); err != nil {
			return err
		}

		return nil
	}
}

func (a *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return a.authenticationWithOptions(next, false)
}

func (a *AuthMiddleware) AuthenticateWithCleanup(next echo.HandlerFunc) echo.HandlerFunc {
	return a.authenticationWithOptions(next, true)
}

func (a *AuthMiddleware) authenticationWithOptions(next echo.HandlerFunc, deleteAfter bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := a.validateClient(c); err != nil {
			return err
		}

		sessionKey := c.Request().Header.Get("X-Session-Key")
		if sessionKey == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Session key required")
		}

		info, err := a.validateSession(sessionKey)
		if err != nil {
			return err
		}

		c.Set(SessionKey, info.ID)

		err = next(c)

		if deleteAfter && err != nil {
			a.deleteSession(sessionKey)
		}
		return err
	}
}
func GetSessionID(c echo.Context) (string, error) {
	value := c.Get(SessionKey)
	if value == nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "invalid(no session)")
	}

	key, ok := value.(string)
	if !ok {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "invalid(invalid format)")
	}

	return key, nil
}

func (a *AuthMiddleware) validateClient(c echo.Context) error {
	clientIP := c.RealIP()
	if a.cfg.ODC.IP != "" && a.cfg.ODC.IP != clientIP {
		log.Printf("Access denied: IP %s is not allowed (expected: %s)", clientIP, a.cfg.ODC.IP)
		return echo.NewHTTPError(http.StatusForbidden, "IP not allowed")
	}

	containerId := c.Request().Header.Get("X-Container-ID")
	if a.cfg.ODC.ID != "" && a.cfg.ODC.ID != containerId {
		log.Printf("Access denied: ID %s is not allowed (expected: %s)", containerId, a.cfg.ODC.ID)
		return echo.NewHTTPError(http.StatusForbidden, "ID not allowed")
	}

	return nil
}

func (a *AuthMiddleware) createSession(c echo.Context) error {
	id := c.Request().Header.Get("X-SEED-ID")
	sessionKey := utils.HashB64([]byte(id))

	a.mutex.Lock()
	defer a.mutex.Unlock()

	if _, exists := a.sessions[sessionKey]; exists {
		return echo.NewHTTPError(http.StatusConflict, "Session already exists")
	}

	a.sessions[sessionKey] = &Session{
		ID:        id,
		CreatedAt: time.Now(),
	}

	return nil
}

func (a *AuthMiddleware) validateSession(sessionKey string) (*Session, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	sessionInfo, exists := a.sessions[sessionKey]

	if !exists {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired session")
	}

	return sessionInfo, nil
}

func (a *AuthMiddleware) deleteSession(sessionKey string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	delete(a.sessions, sessionKey)
}

// ExpiresAt auto clear...
