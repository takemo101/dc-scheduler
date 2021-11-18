package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// SessionAdminAuth admin auth
type SessionAdminAuth struct {
	logger       core.Logger
	app          core.Application
	path         core.Path
	sessionStore core.SessionStore
}

// NewSessionAdminAuth is create middleware
func NewSessionAdminAuth(
	logger core.Logger,
	app core.Application,
	path core.Path,
	sessionStore core.SessionStore,
) SessionAdminAuth {
	return SessionAdminAuth{
		logger:       logger,
		app:          app,
		path:         path,
		sessionStore: sessionStore,
	}
}

// Setup session admin auth middleware
func (m SessionAdminAuth) Setup() {
	m.logger.Info("setup session auth admin")
	m.app.App.Use(m.CreateHandler(true, "system/auth/login"))
}

// CreateHandler is create middleware handler
func (m SessionAdminAuth) CreateHandler(login bool, redirect string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth, err := support.CreateSessionAdminAuthContext(
			c,
			m.sessionStore,
		)
		if err != nil {
			return err
		}

		ok := auth.IsLogin()

		if (login && ok) || (!login && !ok) {
			return c.Next()
		}

		return c.Redirect(m.path.URL(redirect))
	}
}
