package form

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/pkg/domain"
	"github.com/thoas/go-funk"
)

// --- Rules ---

// createBotNameRules 名前ルール
func createBotNameRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("ボット名を入力してください"),
		validation.RuneLength(0, 100).Error("ボット名は80字以内で入力してください"),
	}
}

// createBotAtatarRules アバターファイルルール
func createBotAtatarRules(c *fiber.Ctx, field string) []validation.Rule {
	return []validation.Rule{
		validation.By(func(value interface{}) error {
			file, err := c.FormFile(field)
			if err != nil {
				return nil
			}

			contentType := file.Header.Get("Content-Type")
			checkTypes := domain.GetBotAtatarContentTypes()

			if !funk.Contains(checkTypes, contentType) {
				return errors.New("画像ファイルを選択してください")
			}

			return nil
		}),
	}
}

// createBotWebhookRules ウェブフックURL
func createBotWebhookRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("ウェブフックURLを入力してください"),
		validation.RuneLength(0, 1000).Error("ウェブフックURLは1000字以内で入力してください"),
		validation.Match(regexp.MustCompile("^" + domain.BotDiscordWebhookURLPrefix)).Error("ディスコードのウェブフックURLを入力してください"),
	}
}

// --- BotSearch ---

// BotSearch Bot検索パラメータ
type BotSearch struct {
	Page int `json:"page" form:"page"`
}

// --- BotCreate ---

// BotCreate Admin追加パラメータ
type BotCreateAndUpdate struct {
	Name    string `json:"name" form:"name"`
	Atatar  string `json:"avatar" form:"avatar"`
	Webhook string `json:"webhook" form:"webhook"`
	Active  string `json:"active" form:"active"`
}

// Validate Bot追加バリデーション
func (form BotCreateAndUpdate) Validate(c *fiber.Ctx) error {
	nameRules := createBotNameRules()
	avatarRules := createBotAtatarRules(c, "avatar")
	webhookRules := createBotWebhookRules()

	fields := []*validation.FieldRules{
		validation.Field(
			&form.Name,
			nameRules...,
		),
		validation.Field(
			&form.Atatar,
			avatarRules...,
		),
		validation.Field(
			&form.Webhook,
			webhookRules...,
		),
	}

	return validation.ValidateStruct(&form, fields...)
}

// Sanitize BotCreateの変換
func (form *BotCreateAndUpdate) Sanitize() (err error) {
	form.Name = strings.TrimSpace(form.Name)
	form.Webhook = strings.TrimSpace(form.Webhook)

	return err
}

// ActiveToBool Activeの文字列をboolに変換して返す
func (form BotCreateAndUpdate) ActiveToBool() bool {
	ok, _ := strconv.ParseBool(form.Active)

	return ok
}
