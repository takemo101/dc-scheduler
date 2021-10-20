package form

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
)

// --- Rules ---

// createImmediatePostMessageRules 名前ルール
func createImmediatePostMessageRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("メッセージを入力してください"),
		validation.RuneLength(0, 2000).Error("メッセージは2000字以内で入力してください"),
	}
}

// --- PostMessageSearch ---

// PostMessageSearch PostMessage検索パラメータ
type PostMessageSearch struct {
	Page int `json:"page" form:"page"`
}

// --- Create ---

// ImmediatePostCreate Admin追加パラメータ
type ImmediatePostCreate struct {
	Message string `json:"message" form:"message"`
}

// Validate ImmediatePost追加バリデーション
func (form ImmediatePostCreate) Validate(c *fiber.Ctx) error {
	messageRules := createImmediatePostMessageRules()

	fields := []*validation.FieldRules{
		validation.Field(
			&form.Message,
			messageRules...,
		),
	}

	return validation.ValidateStruct(&form, fields...)
}

// Sanitize ImmediatePostCreateの変換
func (form *ImmediatePostCreate) Sanitize() (err error) {
	form.Message = strings.TrimSpace(form.Message)

	return err
}
