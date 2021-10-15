package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/imdario/mergo"
	"github.com/takemo101/dc-scheduler/app/middleware"
)

type Toastr string

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

type Messages middleware.SessionMessages

func SetToastr(c *fiber.Ctx, ttype Toastr, message string, data Messages) error {
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

	return SetMessages(c, messages)
}

func SetMessages(c *fiber.Ctx, messages Messages) error {
	return middleware.SetSessionMessages(c, middleware.SessionMessages(messages))
}
