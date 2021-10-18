package form

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- Rules ---

// createAdminNameRules 名前ルール
func createAdminNameRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("名前を入力してください"),
		validation.RuneLength(0, 100).Error("名前は100字以内で入力してください"),
	}
}

// createAdminEmailRules メールアドレスルール
func createAdminEmailRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("名前を入力してください"),
		validation.RuneLength(0, 100).Error("名前は100字以内で入力してください"),
		is.Email.Error("メールアドレスは[xxx@xxx.com]のような形式で入力してください"),
	}
}

// createAdminPasswordRules パスワードルール
func createAdminPasswordRules(passwordConfirm string) []validation.Rule {
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

// createAdminPasswordConfirmRules パスワード（確認）ルール
func createAdminPasswordConfirmRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("パスワード（確認）を入力してください"),
	}
}

// createAdminRoleRules ロールルール
func createAdminRoleRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("権限を選択してください"),
		validation.NotIn(
			domain.AdminRoleSystem,
			domain.AdminRoleNormal,
		).Error("権限に正しい値を設定してください"),
	}
}

// --- AdminSearch ---

// AdminSearch Admin検索パラメータ
type AdminSearch struct {
	Page int `json:"page" form:"page"`
}

// --- AdminCreate ---

// AdminCreate Admin追加パラメータ
type AdminCreate struct {
	Name            string `json:"name" form:"name"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm"`
	Role            string `json:"role" form:"role"`
}

// Validate Admin追加バリデーション
func (form AdminCreate) Validate() error {
	nameRules := createAdminNameRules()
	emailRules := createAdminEmailRules()
	passwordRules := createAdminPasswordRules(form.PasswordConfirm)
	passwordConfirmRules := createAdminPasswordConfirmRules()
	roleRules := createAdminRoleRules()

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
		validation.Field(
			&form.Role,
			roleRules...,
		),
	}

	return validation.ValidateStruct(&form, fields...)
}

// Sanitize AdminCreateの変換
func (form *AdminCreate) Sanitize() (err error) {
	form.Name = strings.TrimSpace(form.Name)
	form.Email = strings.TrimSpace(form.Email)
	form.Password = strings.TrimSpace(form.Password)
	form.PasswordConfirm = strings.TrimSpace(form.PasswordConfirm)

	return err
}

// --- AdminUpdate ---

// AdminUpdate Admin更新パラメータ
type AdminUpdate struct {
	Name            string `json:"name" form:"name"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm"`
	Role            string `json:"role" form:"role"`
}

// Validate Admin更新バリデーション
func (form AdminUpdate) Validate() error {
	nameRules := createAdminNameRules()
	emailRules := createAdminEmailRules()
	roleRules := createAdminRoleRules()

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
			&form.Role,
			roleRules...,
		),
	}

	if form.Password != "" {
		passwordRules := createAdminPasswordRules(form.PasswordConfirm)
		passwordConfirmRules := createAdminPasswordConfirmRules()

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

// Sanitize AdminUpdateの変換
func (form *AdminUpdate) Sanitize() (err error) {
	form.Name = strings.TrimSpace(form.Name)
	form.Email = strings.TrimSpace(form.Email)
	form.Password = strings.TrimSpace(form.Password)
	form.PasswordConfirm = strings.TrimSpace(form.PasswordConfirm)

	return err
}

// --- AccountUpdate ---

// AccountUpdate Admin自身アカウント更新パラメータ
type AccountUpdate struct {
	Name            string `json:"name" form:"name"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm"`
}

// Validate Admin自身アカウント更新バリデーション
func (form AccountUpdate) Validate() error {
	nameRules := createAdminNameRules()
	emailRules := createAdminEmailRules()

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
		passwordRules := createAdminPasswordRules(form.PasswordConfirm)
		passwordConfirmRules := createAdminPasswordConfirmRules()

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

// Sanitize AdminUpdateの変換
func (form *AccountUpdate) Sanitize() (err error) {
	form.Name = strings.TrimSpace(form.Name)
	form.Email = strings.TrimSpace(form.Email)
	form.Password = strings.TrimSpace(form.Password)
	form.PasswordConfirm = strings.TrimSpace(form.PasswordConfirm)

	return err
}
