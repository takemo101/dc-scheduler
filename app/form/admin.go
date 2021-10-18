package form

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

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
	fields := []*validation.FieldRules{
		validation.Field(
			&form.Name,
			validation.Required.Error("名前を入力してください"),
			validation.RuneLength(0, 100).Error("名前は100字以内で入力してください"),
		),
		validation.Field(
			&form.Email,
			validation.Required.Error("メールアドレスを入力してください"),
			validation.Length(0, 180).Error("メールアドレスは180字以内で入力してください"),
			is.Email.Error("メールアドレスは[xxx@xxx.com]のような形式で入力してください"),
		),
		validation.Field(
			&form.Password,
			validation.Required.Error("パスワードを入力してください"),
			validation.Length(3, 20).Error("パスワードは3から20文字以内で入力してください"),
			validation.By(func(value interface{}) error {
				if form.PasswordConfirm != value.(string) {
					return errors.New("パスワードが一致しません")
				}
				return nil
			}),
		),
		validation.Field(
			&form.PasswordConfirm,
			validation.Required.Error("パスワード（確認）を入力してください"),
		),
		validation.Field(
			&form.Role,
			validation.Required.Error("権限を選択してください"),
			validation.NotIn(
				domain.AdminRoleSystem,
				domain.AdminRoleNormal,
			).Error("権限に正しい値を設定してください"),
		),
	}

	return validation.ValidateStruct(&form, fields...)
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
	fields := []*validation.FieldRules{
		validation.Field(
			&form.Name,
			validation.Required.Error("名前を入力してください"),
			validation.RuneLength(0, 100).Error("名前は100字以内で入力してください"),
		),
		validation.Field(
			&form.Email,
			validation.Required.Error("メールアドレスを入力してください"),
			validation.Length(0, 180).Error("メールアドレスは180字以内で入力してください"),
			is.Email.Error("メールアドレスは[xxx@xxx.com]のような形式で入力してください"),
		),
		validation.Field(
			&form.Role,
			validation.Required.Error("権限を選択してください"),
			validation.NotIn(
				domain.AdminRoleSystem,
				domain.AdminRoleNormal,
			).Error("権限に正しい値を設定してください"),
		),
	}

	if form.Password != "" {
		fields = append(
			fields,
			validation.Field(
				&form.Password,
				validation.Required.Error("パスワードを入力してください"),
				validation.Length(3, 20).Error("パスワードは3から20文字以内で入力してください"),
				validation.By(func(value interface{}) error {
					if form.PasswordConfirm != value.(string) {
						return errors.New("パスワードが一致しません")
					}
					return nil
				}),
			),
			validation.Field(
				&form.PasswordConfirm,
				validation.Required.Error("パスワード（確認）を入力してください"),
			),
		)
	}

	return validation.ValidateStruct(&form, fields...)
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
	fields := []*validation.FieldRules{
		validation.Field(
			&form.Name,
			validation.Required.Error("名前を入力してください"),
			validation.RuneLength(0, 100).Error("名前は100字以内で入力してください"),
		),
		validation.Field(
			&form.Email,
			validation.Required.Error("メールアドレスを入力してください"),
			validation.Length(0, 180).Error("メールアドレスは180字以内で入力してください"),
			is.Email.Error("メールアドレスは[xxx@xxx.com]のような形式で入力してください"),
		),
	}

	if form.Password != "" {
		fields = append(
			fields,
			validation.Field(
				&form.Password,
				validation.Required.Error("パスワードを入力してください"),
				validation.Length(3, 20).Error("パスワードは3から20文字以内で入力してください"),
				validation.By(func(value interface{}) error {
					if form.PasswordConfirm != value.(string) {
						return errors.New("パスワードが一致しません")
					}
					return nil
				}),
			),
			validation.Field(
				&form.PasswordConfirm,
				validation.Required.Error("パスワード（確認）を入力してください"),
			),
		)
	}

	return validation.ValidateStruct(&form, fields...)
}
