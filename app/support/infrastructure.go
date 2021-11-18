package support

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/domain"
	"github.com/takemo101/dc-scheduler/pkg/infrastructure"
)

// CreateSessionAdminAuthContext セッションによるAdminAuthContextをinfrastructureから生成
func CreateSessionAdminAuthContext(
	c *fiber.Ctx,
	sessionStore core.SessionStore,
) (auth domain.AdminAuthContext, err error) {
	session, err := sessionStore.GetSession(c)
	if err != nil {
		return auth, err
	}

	return infrastructure.NewAdminAuthContext(session), err
}

// CreateSessionUserAuthContext セッションによるUserAuthContextをinfrastructureから生成
func CreateSessionUserAuthContext(
	c *fiber.Ctx,
	sessionStore core.SessionStore,
) (auth domain.UserAuthContext, err error) {
	session, err := sessionStore.GetSession(c)
	if err != nil {
		return auth, err
	}

	return infrastructure.NewUserAuthContext(session), err
}
