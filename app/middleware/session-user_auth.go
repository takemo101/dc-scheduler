package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// SessionUserAuth user auth
type SessionUserAuth struct {
	logger       core.Logger
	app          core.Application
	path         core.Path
	sessionStore core.SessionStore
}

// NewSessionUserAuth is create middleware
func NewSessionUserAuth(
	logger core.Logger,
	app core.Application,
	path core.Path,
	sessionStore core.SessionStore,
) SessionUserAuth {
	return SessionUserAuth{
		logger:       logger,
		app:          app,
		path:         path,
		sessionStore: sessionStore,
	}
}

// Setup session user auth middleware
func (m SessionUserAuth) Setup() {
	m.logger.Info("setup session auth user")
	m.app.App.Use(m.CreateHandler(true, "auth/login"))
}

// CreateHandler is create middleware handler
func (m SessionUserAuth) CreateHandler(login bool, redirect string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth, err := support.CreateSessionUserAuthContext(
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
