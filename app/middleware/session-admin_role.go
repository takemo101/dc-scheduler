package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// SessionAdminRole admin auth
type SessionAdminRole struct {
	logger       core.Logger
	value        support.ContextValue
	sessionStore core.SessionStore
}

// NewCheckAdminRole is create middleware
func NewSessionAdminRole(
	logger core.Logger,
	value support.ContextValue,
	sessionStore core.SessionStore,
) SessionAdminRole {
	return SessionAdminRole{
		logger:       logger,
		sessionStore: sessionStore,
	}
}

// CreateHandler is create middleware handler
func (m SessionAdminRole) CreateHandler(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth, err := support.CreateSessionAdminAuthContext(
			c,
			m.sessionStore,
		)
		if err != nil {
			return err
		}

		admin, err := auth.AdminAuth()
		if err != nil {
			return err
		}

		for _, role := range roles {
			if admin.HaveRole(role) {
				return c.Next()
			}
		}

		response := m.value.GetResponseHelper(c)
		return response.Error(fiber.ErrUnauthorized, fiber.StatusUnauthorized)
	}
}
