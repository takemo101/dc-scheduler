package form

import (
	"errors"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// --- Rules ---

// createUserNameRules 名前ルール
func createUserNameRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("名前を入力してください"),
		validation.RuneLength(0, 100).Error("名前は100字以内で入力してください"),
	}
}

// createUserEmailRules メールアドレスルール
func createUserEmailRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("名前を入力してください"),
		validation.RuneLength(0, 100).Error("名前は100字以内で入力してください"),
		is.Email.Error("メールアドレスは[xxx@xxx.com]のような形式で入力してください"),
	}
}

// createUserPasswordRules パスワードルール
func createUserPasswordRules(passwordConfirm string) []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("パスワードを入力してください"),
		validation.Length(3, 20).Error("パスワードは3から20文字以内で入力してください"),
		validation.By(func(value interface{}) error {
			if passwordConfirm != value.(string) {
				return errors.New("パスワードが一致しません")
			}
			return nil
		}),
	}
}

// createUserPasswordConfirmRules パスワード（確認）ルール
func createUserPasswordConfirmRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("パスワード（確認）を入力してください"),
	}
}

// --- UserSearch ---

// UserSearch User検索パラメータ
type UserSearch struct {
	Page int `json:"page" form:"page"`
}

// --- UserCreate ---

// UserCreate User追加パラメータ
type UserCreate struct {
	Name            string `json:"name" form:"name"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm"`
	Active          string `json:"active" form:"active"`
}

// Validate User追加バリデーション
func (form UserCreate) Validate() error {
	nameRules := createUserNameRules()
	emailRules := createUserEmailRules()
	passwordRules := createUserPasswordRules(form.PasswordConfirm)
	passwordConfirmRules := createUserPasswordConfirmRules()

	fields := []*validation.FieldRules{
		validation.Field(
			&form.Name,
			nameRules...,
		),
		validation.Field(
			&form.Email,
			emailRules...,
		),
		validation.Field(
			&form.Password,
			passwordRules...,
		),
		validation.Field(
			&form.PasswordConfirm,
			passwordConfirmRules...,
		),
	}

	return validation.ValidateStruct(&form, fields...)
}

// Sanitize UserCreateの変換
func (form *UserCreate) Sanitize() (err error) {
	form.Name = strings.TrimSpace(form.Name)
	form.Email = strings.TrimSpace(form.Email)
	form.Password = strings.TrimSpace(form.Password)
	form.PasswordConfirm = strings.TrimSpace(form.PasswordConfirm)

	return err
}

// ActiveToBool Activeの文字列をboolに変換して返す
func (form UserCreate) ActiveToBool() bool {
	ok, _ := strconv.ParseBool(form.Active)

	return ok
}

// --- UserUpdate ---

// UserUpdate User更新パラメータ
type UserUpdate struct {
	Name            string `json:"name" form:"name"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm"`
	Active          string `json:"active" form:"active"`
}

// Validate User更新バリデーション
func (form UserUpdate) Validate() error {
	nameRules := createUserNameRules()
	emailRules := createUserEmailRules()

	fields := []*validation.FieldRules{
		validation.Field(
			&form.Name,
			nameRules...,
		),
		validation.Field(
			&form.Email,
			emailRules...,
		),
	}

	if form.Password != "" {
		passwordRules := createUserPasswordRules(form.PasswordConfirm)
		passwordConfirmRules := createUserPasswordConfirmRules()

		fields = append(
			fields,
			validation.Field(
				&form.Password,
				passwordRules...,
			),
			validation.Field(
				&form.PasswordConfirm,
				passwordConfirmRules...,
			),
		)
	}

	return validation.ValidateStruct(&form, fields...)
}

// Sanitize UserUpdateの変換
func (form *UserUpdate) Sanitize() (err error) {
	form.Name = strings.TrimSpace(form.Name)
	form.Email = strings.TrimSpace(form.Email)
	form.Password = strings.TrimSpace(form.Password)
	form.PasswordConfirm = strings.TrimSpace(form.PasswordConfirm)

	return err
}

// ActiveToBool Activeの文字列をboolに変換して返す
func (form UserUpdate) ActiveToBool() bool {
	ok, _ := strconv.ParseBool(form.Active)

	return ok
}
