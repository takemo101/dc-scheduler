package core

import (
	"encoding/gob"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
)

// SessionStore is fibar session store
type SessionStore struct {
	store *session.Store
}

func NewSessionStore(
	config Config,
) SessionStore {

	gob.Register(SessionErrors{})
	gob.Register(SessionInputs{})
	gob.Register(SessionMessages{})

	store := session.New(session.Config{
		Expiration:     config.Session.Expiration,
		KeyLookup:      "cookie:" + config.Session.Name,
		CookieDomain:   config.Session.Domain,
		CookiePath:     config.Session.Path,
		CookieSecure:   config.Session.Secure,
		CookieHTTPOnly: config.Session.HTTPOnly,
		KeyGenerator:   utils.UUID,
	})

	return SessionStore{
		store,
	}
}

// GetSession is fiber context to get session
func (ss *SessionStore) GetSession(c *fiber.Ctx) (*session.Session, error) {
	return ss.store.Get(c)
}

type (
	SessionErrors   map[string]string
	SessionInputs   map[string]interface{}
	SessionMessages map[string]interface{}
)

const (
	SessionErrorsKey   string = "session-errors"
	SessionInputsKey   string = "session-inputs"
	SessionMessagesKey string = "session-messages"
)

// GetSessionErrors session errors process
func (ss *SessionStore) GetSessionErrors(c *fiber.Ctx) (SessionErrors, error) {
	sessionErrors, err := ss.GetSessionValue(c, SessionErrorsKey)
	if err != nil {
		return nil, err
	}
	if sessionErrors != nil {
		return sessionErrors.(SessionErrors), nil
	}
	return nil, errors.New("not found errors")
}

// SetSessionErrors session errors process
func (ss *SessionStore) SetSessionErrors(c *fiber.Ctx, errors SessionErrors) error {
	return ss.SetSessionValue(c, SessionErrorsKey, errors)
}

// GetSessionInputs session old inputs process
func (ss *SessionStore) GetSessionInputs(c *fiber.Ctx) (SessionInputs, error) {
	inputs, err := ss.GetSessionValue(c, SessionInputsKey)
	if err != nil {
		return nil, err
	}
	if inputs != nil {
		return inputs.(SessionInputs), nil
	}
	return nil, errors.New("not found inputs")
}

// SetSessionInputs session old inputs process
func (ss *SessionStore) SetSessionInputs(c *fiber.Ctx, inputs SessionInputs) error {
	return ss.SetSessionValue(c, SessionInputsKey, inputs)
}

// GetSessionValue get session value
func (ss *SessionStore) GetSessionValue(c *fiber.Ctx, key string) (interface{}, error) {
	session, err := ss.GetSession(c)
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
func (ss *SessionStore) SetSessionValue(c *fiber.Ctx, key string, value interface{}) error {
	session, err := ss.GetSession(c)
	if err != nil {
		return err
	}
	defer session.Save()
	session.Set(key, value)
	return nil
}

// SetErrorResource set session inputs and errors
func (ss *SessionStore) SetErrorResource(c *fiber.Ctx, errors SessionErrors, inputs SessionInputs) {
	ss.SetSessionErrors(c, errors)
	ss.SetSessionInputs(c, inputs)
}

// GetSessionMessages session flash messages process
func (ss *SessionStore) GetSessionMessages(c *fiber.Ctx) (SessionMessages, error) {
	messages, err := ss.GetSessionValue(c, SessionMessagesKey)
	if err != nil {
		return nil, err
	}
	if messages != nil {
		return messages.(SessionMessages), nil
	}
	return nil, errors.New("not found messages")
}

// SetSessionMessages session flash messages process
func (ss *SessionStore) SetSessionMessages(c *fiber.Ctx, messages SessionMessages) error {
	return ss.SetSessionValue(c, SessionMessagesKey, messages)
}
