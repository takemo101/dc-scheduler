package support

import (
	"github.com/gofiber/fiber/v2"
	"github.com/imdario/mergo"
	"github.com/takemo101/dc-scheduler/core"
)

type ToastrMessage struct {
	sessionStore core.SessionStore
}

func NewToastrMessage(
	sessionStore core.SessionStore,
) ToastrMessage {
	return ToastrMessage{
		sessionStore,
	}
}

type (
	Toastr   string
	Messages core.SessionMessages
)

const (
	ToastrError   Toastr = "error"
	ToastrSuccess Toastr = "success"
	ToastrStore   Toastr = "store"
	ToastrUpdate  Toastr = "update"
	ToastrDelete  Toastr = "delete"
)

func (t Toastr) String() string {
	return string(t)
}

func (t Toastr) Message() string {
	switch t {
	case ToastrError:
		return "入力内容を確認してください"
	case ToastrStore:
		return "追加しました"
	case ToastrUpdate:
		return "更新しました"
	case ToastrDelete:
		return "削除しました"
	}
	return ""
}

func (tm *ToastrMessage) SetToastr(
	c *fiber.Ctx,
	ttype Toastr,
	message string,
	data Messages,
) error {
	messages := Messages{
		"toastr_type": string(ttype),
		"message":     message,
	}
	err := mergo.Merge(
		&messages,
		data,
	)

	if err != nil {
		return err
	}

	return tm.SetMessages(c, messages)
}

func (tm *ToastrMessage) SetMessages(
	c *fiber.Ctx,
	messages Messages,
) error {
	return tm.sessionStore.SetSessionMessages(
		c,
		core.SessionMessages(messages),
	)
}
