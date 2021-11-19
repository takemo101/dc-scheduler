package form

import (
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/takemo101/dc-scheduler/pkg/domain"
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

// SchedulePostCreateAndUpdate SchedulePost追加・パラメータ
type SchedulePostCreateAndUpdate struct {
	Message       string `json:"message" form:"message"`
	ReservationAt string `json:"reservation_at" form:"reservation_at"`
}

// Validate SchedulePost追加バリデーション
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
	at, _ := time.ParseInLocation("2006-01-02 15:04", form.ReservationAt, time.Local)
	return at
}

// --- RegularPostCreateAndUpdate ---

// RegularPostCreateAndUpdate RegularPost追加・更新パラメータ
type RegularPostCreateAndUpdate struct {
	Message string `json:"message" form:"message"`
	Active  string `json:"active" form:"active"`
}

// Validate RegularPost追加・更新バリデーション
func (form RegularPostCreateAndUpdate) Validate() error {
	messageRules := createPostMessageMessageRules()

	fields := []*validation.FieldRules{
		validation.Field(
			&form.Message,
			messageRules...,
		),
	}

	return validation.ValidateStruct(&form, fields...)
}

// Sanitize RegularPostCreateAndUpdateの変換
func (form *RegularPostCreateAndUpdate) Sanitize() (err error) {
	form.Message = strings.TrimSpace(form.Message)

	return err
}

// ActiveToBool Activeの文字列をboolに変換して返す
func (form RegularPostCreateAndUpdate) ActiveToBool() bool {
	ok, _ := strconv.ParseBool(form.Active)

	return ok
}

// --- RegularTimingAdd ---

// RegularTimingAdd RegularTiming追加パラメータ
type RegularTimingAdd struct {
	DayOfWeek string `json:"day_of_week" form:"day_of_week"`
	HourTime  string `json:"hour_time" form:"hour_time"`
}

// Validate RegularTiming追加バリデーション
func (form RegularTimingAdd) Validate() error {
	fields := []*validation.FieldRules{
		validation.Field(
			&form.DayOfWeek,
			validation.Required.Error("配信曜日を選択してください"),
			validation.NotIn(
				domain.DayOfWeekSunday,
				domain.DayOfWeekMonday,
				domain.DayOfWeekTuesday,
				domain.DayOfWeekWednesday,
				domain.DayOfWeekThursday,
				domain.DayOfWeekFriday,
				domain.DayOfWeekSaturday,
			).Error("配信曜日に正しい値を設定してください"),
		),
		validation.Field(
			&form.HourTime,
			validation.Required.Error("配信時間を入力してください"),
			validation.Date("15:04").Error("配信時間の形式が不正です"),
		),
	}

	return validation.ValidateStruct(&form, fields...)
}

// HourTimeToTime RegularTimingAddの文字列をtime.Timeに変換して返す
func (form RegularTimingAdd) HourTimeToTime() time.Time {
	at, _ := time.ParseInLocation("15:04", form.HourTime, time.Local)
	return at
}

// --- ApiPostSend ---

// ApiPostSend ApiPost送信パラメータ
type ApiPostSend struct {
	Message string `json:"message" form:"message"`
}

// Validate ApiPost送信パラメータ
func (form ApiPostSend) Validate() error {
	messageRules := createPostMessageMessageRules()

	fields := []*validation.FieldRules{
		validation.Field(
			&form.Message,
			messageRules...,
		),
	}

	return validation.ValidateStruct(&form, fields...)
}

// Sanitize ApiPostSendの変換
func (form *ApiPostSend) Sanitize() (err error) {
	form.Message = strings.TrimSpace(form.Message)

	return err
}
