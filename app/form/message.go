package form

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// --- Rules ---

// createPostMessageMessageRules メッセージルール
func createPostMessageMessageRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("メッセージを入力してください"),
		validation.RuneLength(0, 2000).Error("メッセージは2000字以内で入力してください"),
	}
}

// createSchedulePostReservationAtRules 名前ルール
func createSchedulePostReservationAtRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("予約日時を設定してください"),
		validation.Date("2006-01-02 15:04").Error("予約日時の形式が不正です").Min(time.Now()).RangeError("予約日時は現在以降を指定してください"),
	}
}

// --- PostMessageSearch ---

// PostMessageSearch PostMessage検索パラメータ
type PostMessageSearch struct {
	Page int `json:"page" form:"page"`
}

// --- SentMessageHistory ---

// SentMessageHistory SentMessage検索パラメータ
type SentMessageHistory struct {
	Page int `json:"page" form:"page"`
}

// --- ImmediatePostCreate ---

// ImmediatePostCreate Admin追加パラメータ
type ImmediatePostCreate struct {
	Message string `json:"message" form:"message"`
}

// Validate ImmediatePost追加バリデーション
func (form ImmediatePostCreate) Validate() error {
	messageRules := createPostMessageMessageRules()

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

// --- SchedulePostCreateAndUpdate ---

// SchedulePostCreateAndUpdate Admin追加・パラメータ
type SchedulePostCreateAndUpdate struct {
	Message       string `json:"message" form:"message"`
	ReservationAt string `json:"reservation_at" form:"reservation_at"`
}

// Validate ImmediatePost追加バリデーション
func (form SchedulePostCreateAndUpdate) Validate() error {
	messageRules := createPostMessageMessageRules()
	reservationAtRules := createSchedulePostReservationAtRules()

	fields := []*validation.FieldRules{
		validation.Field(
			&form.Message,
			messageRules...,
		),
		validation.Field(
			&form.ReservationAt,
			reservationAtRules...,
		),
	}

	return validation.ValidateStruct(&form, fields...)
}

// Sanitize SchedulePostCreateAndUpdateの変換
func (form *SchedulePostCreateAndUpdate) Sanitize() (err error) {
	form.Message = strings.TrimSpace(form.Message)

	return err
}

// ReservationAtToTime ReservationAtの文字列をtime.Timeに変換して返す
func (form SchedulePostCreateAndUpdate) ReservationAtToTime() time.Time {
	at, _ := time.Parse("2006-01-02 15:04", form.ReservationAt)

	return at
}
