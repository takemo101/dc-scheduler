package form

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Login ログインフォームパラメータ
type Login struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// Validate ログインバリデーション
func (form Login) Validate() error {
	fields := []*validation.FieldRules{
		validation.Field(
			&form.Email,
			validation.Required.Error("メールアドレスを入力してください"),
		),
		validation.Field(
			&form.Password,
			validation.Required.Error("パスワードを入力してください"),
		),
	}

	return validation.ValidateStruct(&form, fields...)
}
