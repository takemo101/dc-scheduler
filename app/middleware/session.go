package middleware

import (
	"encoding/gob"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/takemo101/dc-scheduler/core"
)

// Session is struct
type Session struct {
	logger core.Logger
	app    core.Application
	config core.Config
}

// NewSession is create middleware
func NewSession(
	logger core.Logger,
	app core.Application,
	config core.Config,
) Session {
	return Session{
		logger: logger,
		app:    app,
		config: config,
	}
}

// Setup start session middleware
func (m Session) Setup() {
	m.logger.Info("setup session")

	gob.Register(SessionErrors{})
	gob.Register(SessionInputs{})
	gob.Register(SessionMessages{})

	store := session.New(session.Config{
		Expiration:     m.config.Session.Expiration,
		KeyLookup:      "cookie:" + m.config.Session.Name,
		CookieDomain:   m.config.Session.Domain,
		CookiePath:     m.config.Session.Path,
		CookieSecure:   m.config.Session.Secure,
		CookieHTTPOnly: m.config.Session.HTTPOnly,
		KeyGenerator:   utils.UUID,
	})

	m.app.App.Use(func(c *fiber.Ctx) error {
		c.Locals("session", store)
		return c.Next()
	})
}

// GetSession to session
func GetSession(c *fiber.Ctx) (*session.Session, error) {
	store := c.Locals("session").(*session.Store)
	return store.Get(c)
}

type SessionErrors map[string]string

// GetSessionErrors session errors process
func GetSessionErrors(c *fiber.Ctx) (SessionErrors, error) {
	sessionErrors, err := GetSessionValue(c, "session-errors")
	if err != nil {
		return nil, err
	}
	if sessionErrors != nil {
		return sessionErrors.(SessionErrors), nil
	}
	return nil, errors.New("not found errors")
}

// SetSessionErrors session errors process
func SetSessionErrors(c *fiber.Ctx, errors SessionErrors) error {
	return SetSessionValue(c, "session-errors", errors)
}

type SessionInputs map[string]interface{}

// GetSessionInputs session old inputs process
func GetSessionInputs(c *fiber.Ctx) (SessionInputs, error) {
	inputs, err := GetSessionValue(c, "session-inputs")
	if err != nil {
		return nil, err
	}
	if inputs != nil {
		return inputs.(SessionInputs), nil
	}
	return nil, errors.New("not found inputs")
}

// SetSessionInputs session old inputs process
func SetSessionInputs(c *fiber.Ctx, inputs SessionInputs) error {
	return SetSessionValue(c, "session-inputs", inputs)
}

// GetSessionValue get session value
func GetSessionValue(c *fiber.Ctx, key string) (interface{}, error) {
	session, err := GetSession(c)
	if err != nil {
		return nil, err
	}

	if value := session.Get(key); value != nil {
		defer session.Save()
		session.Set(key, nil)
		return value, nil
	}
	return nil, nil
}

// SetSessionValue set session value
func SetSessionValue(c *fiber.Ctx, key string, value interface{}) error {
	session, err := GetSession(c)
	if err != nil {
		return err
	}
	defer session.Save()
	session.Set(key, value)
	return nil
}

// SetErrorResource set session inputs and errors
func SetErrorResource(c *fiber.Ctx, errors SessionErrors, inputs SessionInputs) {
	SetSessionErrors(c, errors)
	SetSessionInputs(c, inputs)
}

type SessionMessages map[string]interface{}

// GetSessionMessages session flash messages process
func GetSessionMessages(c *fiber.Ctx) (SessionMessages, error) {
	messages, err := GetSessionValue(c, "session-messages")
	if err != nil {
		return nil, err
	}
	if messages != nil {
		return messages.(SessionMessages), nil
	}
	return nil, errors.New("not found messages")
}

// SetSessionMessages session flash messages process
func SetSessionMessages(c *fiber.Ctx, messages SessionMessages) error {
	return SetSessionValue(c, "session-messages", messages)
}
