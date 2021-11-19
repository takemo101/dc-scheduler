package domain

import (
	"errors"
	"time"
)

// --- DayOfWeek ValueObject ---

// DayOfWeek RegularTimingの週設定
type DayOfWeek string

const (
	DayOfWeekSunday    DayOfWeek = "sun"
	DayOfWeekMonday    DayOfWeek = "mon"
	DayOfWeekTuesday   DayOfWeek = "tue"
	DayOfWeekWednesday DayOfWeek = "wed"
	DayOfWeekThursday  DayOfWeek = "thu"
	DayOfWeekFriday    DayOfWeek = "fri"
	DayOfWeekSaturday  DayOfWeek = "sat"
)

// NewDayOfWeek コンストラクタ
func NewDayOfWeek(week string) (vo DayOfWeek, err error) {

	vo = DayOfWeek(week)

	if !vo.Valid() {
		return vo, errors.New("DayOfWeekの値が不正です")
	}

	return vo, err
}

// Value 値を返す
func (vo DayOfWeek) Value() string {
	return vo.String()
}

// String stringへの変換
func (vo DayOfWeek) String() string {
	return string(vo)
}

// Name 日本語名を返す
func (vo DayOfWeek) Name() string {
	switch vo {
	case DayOfWeekSunday:
		return "日"
	case DayOfWeekMonday:
		return "月"
	case DayOfWeekTuesday:
		return "火"
	case DayOfWeekWednesday:
		return "水"
	case DayOfWeekThursday:
		return "木"
	case DayOfWeekFriday:
		return "金"
	case DayOfWeekSaturday:
		return "土"
	}
	return ""
}

// ToWeekday time.Weekdayに変換
func (vo DayOfWeek) ToWeekday() time.Weekday {
	switch vo {
	case DayOfWeekSunday:
		return time.Sunday
	case DayOfWeekMonday:
		return time.Monday
	case DayOfWeekTuesday:
		return time.Tuesday
	case DayOfWeekWednesday:
		return time.Wednesday
	case DayOfWeekThursday:
		return time.Thursday
	case DayOfWeekFriday:
		return time.Friday
	case DayOfWeekSaturday:
		return time.Saturday
	}
	return time.Sunday
}

// Valid 定義したものに一致するか
func (vo DayOfWeek) Valid() bool {
	switch vo {
	case DayOfWeekSunday, DayOfWeekMonday, DayOfWeekTuesday, DayOfWeekWednesday, DayOfWeekThursday, DayOfWeekFriday, DayOfWeekSaturday:
		return true
	}
	return false
}

// Equals VOの値が一致するか
func (vo DayOfWeek) Equals(eq DayOfWeek) bool {
	return vo.Value() == eq.Value()
}

// FromWeekday time.WeekdayからDayOfWeekに変換
func FromWeekday(week time.Weekday) DayOfWeek {
	switch week {
	case time.Sunday:
		return DayOfWeekSunday
	case time.Monday:
		return DayOfWeekMonday
	case time.Tuesday:
		return DayOfWeekTuesday
	case time.Wednesday:
		return DayOfWeekWednesday
	case time.Thursday:
		return DayOfWeekThursday
	case time.Friday:
		return DayOfWeekFriday
	case time.Saturday:
		return DayOfWeekSaturday
	}
	return DayOfWeekSunday
}

// DayOfWeeks 週を配列で返す
func DayOfWeeks() []DayOfWeek {
	return []DayOfWeek{
		DayOfWeekSunday,
		DayOfWeekMonday,
		DayOfWeekTuesday,
		DayOfWeekWednesday,
		DayOfWeekThursday,
		DayOfWeekFriday,
		DayOfWeekSaturday,
	}
}

// DayOfWeekToArray 週をキー値形式で返す
func DayOfWeekToArray() []KeyValue {
	return []KeyValue{
		{
			Key:   string(DayOfWeekSunday),
			Value: DayOfWeekSunday.Name(),
		},
		{
			Key:   string(DayOfWeekMonday),
			Value: DayOfWeekMonday.Name(),
		},
		{
			Key:   string(DayOfWeekTuesday),
			Value: DayOfWeekTuesday.Name(),
		},
		{
			Key:   string(DayOfWeekWednesday),
			Value: DayOfWeekWednesday.Name(),
		},
		{
			Key:   string(DayOfWeekThursday),
			Value: DayOfWeekThursday.Name(),
		},
		{
			Key:   string(DayOfWeekFriday),
			Value: DayOfWeekFriday.Name(),
		},
		{
			Key:   string(DayOfWeekSaturday),
			Value: DayOfWeekSaturday.Name(),
		},
	}
}

// --- HourTime ValueObject ---

// HourTime RegularTimingの時分設定
type HourTime struct {
	at time.Time
}

// NewHourTime コンストラクタ
func NewHourTime(at time.Time) HourTime {
	return HourTime{
		at: time.Date(
			2020,
			10,
			10,
			at.Hour(),
			TimeToIntervalMinute(at),
			0,
			0,
			at.Location(),
		),
	}
}

// Value 値を返す
func (vo HourTime) Value() time.Time {
	return vo.at
}

// String stringへの変換
func (vo HourTime) String() string {
	return vo.Value().Format("15:04")
}

// Equals VOの値が一致するか
func (vo HourTime) Equals(eq HourTime) bool {
	return vo.String() == eq.String()
}

// --- RegularTiming Entity ---'

// RegularTiming メッセージ配信Entity
type RegularTiming struct {
	id        UUID
	dayOfWeek DayOfWeek
	hourTime  HourTime
}

// NewRegularTiming コンストラクタ
func NewRegularTiming(
	week DayOfWeek,
	hourTime time.Time,
) RegularTiming {
	return RegularTiming{
		id:        GenerateUUID(),
		dayOfWeek: week,
		hourTime:  NewHourTime(hourTime),
	}
}

// CreateRegularTiming 定期配信のタイミングを作成
func CreateRegularTiming(
	week string,
	hourTime time.Time,
) (entity RegularTiming, err error) {

	dayOfWeekVO, err := NewDayOfWeek(week)
	if err != nil {
		return entity, err
	}

	return NewRegularTiming(
		dayOfWeekVO,
		hourTime,
	), err
}

// CreateRegularTimingByTime time.Timeから定期配信のタイミングを作成
func CreateRegularTimingByTime(at time.Time) RegularTiming {
	dayOfWeekVO := FromWeekday(at.Weekday())

	return NewRegularTiming(
		dayOfWeekVO,
		at,
	)
}

func (entity RegularTiming) ID() UUID {
	return entity.id
}

func (entity RegularTiming) DayOfWeek() DayOfWeek {
	return entity.dayOfWeek
}

func (entity RegularTiming) HourTime() HourTime {
	return entity.hourTime
}

// Equals RegularTimingが一致するか
func (entity RegularTiming) Equals(eq RegularTiming) bool {
	return entity.HourTime().Equals(eq.HourTime()) && entity.DayOfWeek().Equals(eq.DayOfWeek())
}

// --- RegularPost Entity ---

// RegularPost 定期配信メッセージEntity
type RegularPost struct {
	PostMessage
	active  bool
	timings []RegularTiming
}

// NewRegularPost コンストラクタ
func NewRegularPost(
	id uint,
	message string,
	bot Bot,
	sentMessages []SentMessage,
	active bool,
	timings []RegularTiming,
) RegularPost {
	return RegularPost{
		PostMessage: NewPostMessage(
			id,
			message,
			MessageTypeRegularPost,
			bot,
			sentMessages,
		),
		active:  active,
		timings: timings,
	}
}

// CreateRegularPost 定期配信メッセージを生成する
func CreateRegularPost(
	id uint,
	message string,
	bot Bot,
	active bool,
) (entity RegularPost, err error) {
	postMessage, err := CreatePostMessage(
		id,
		message,
		MessageTypeRegularPost,
		bot,
	)
	if err != nil {
		return entity, err
	}

	return RegularPost{
		PostMessage: postMessage,
		active:      active,
		timings:     []RegularTiming{},
	}, err
}

// Update RegularPostを更新
func (entity *RegularPost) Update(
	message string,
	active bool,
) (err error) {
	entity.active = active

	return entity.ChangeMessage(message)
}

// IsActive RegularPostが有効か
func (entity RegularPost) IsActive() bool {
	return entity.active
}

func (entity RegularPost) Timings() []RegularTiming {
	return entity.timings
}

// AddTiming 配信タイミングの追加
func (entity *RegularPost) AddTiming(week string, hourTime time.Time) (err error) {

	timing, err := CreateRegularTiming(
		week,
		hourTime,
	)
	if err != nil {
		return err
	}

	return entity.AddTimingByEntity(timing)
}

// AddTimingByEntity 配信タイミングをRegularTimingEntityで追加
func (entity *RegularPost) AddTimingByEntity(timing RegularTiming) (err error) {

	for _, tm := range entity.Timings() {
		if timing.Equals(tm) {
			return errors.New("RegularTimingが重複しています")
		}
	}

	entity.timings = append(
		entity.timings,
		timing,
	)

	return err
}

// RemoveTiming 配信タイミングを削除する
func (entity *RegularPost) RemoveTiming(week string, hourTime time.Time) (err error) {

	timing, err := CreateRegularTiming(
		week,
		hourTime,
	)
	if err != nil {
		return err
	}

	return entity.RemoveTimingByEntity(timing)
}

// RemoveTimingByEntity 配信タイミングをRegularTimingEntityで削除
func (entity *RegularPost) RemoveTimingByEntity(timing RegularTiming) (err error) {

	for i, tm := range entity.Timings() {
		if timing.Equals(tm) {
			entity.timings = append(
				entity.timings[:i],
				entity.timings[i+1:]...,
			)

			return err
		}
	}

	return errors.New("対象のRegularTimingがありません")
}

// IsInTiming 指定した時間が定期配信タイミング内か
func (entity RegularPost) IsInTiming(now time.Time) bool {
	timing := CreateRegularTimingByTime(now)

	for _, tm := range entity.Timings() {
		if timing.Equals(tm) {
			return true
		}
	}
	return false
}

// CanSent 定期配信可能か
func (entity RegularPost) CanSent(now time.Time) bool {
	return entity.IsInTiming(now) && entity.IsActive() && entity.PostMessage.CanSent()
}

// Send メッセージを配信する
func (entity *RegularPost) Send(now time.Time) (send SentMessage, err error) {
	if !entity.CanSent(now) {
		return send, errors.New("Message配信可能ではありません")
	}

	send, err = SendMessage(entity.message, now)

	entity.sentMessages = append(
		entity.sentMessages,
		send,
	)

	return send, err
}

// --- RegularPostRepository ---

// RegularPostRepository 定期配信Entityの永続化
type RegularPostRepository interface {
	PostMessageRepository
	SendList(timing RegularTiming) ([]RegularPost, error)
	Store(entity RegularPost) (PostMessageID, error)
	Update(entity RegularPost) error
	FindByID(id PostMessageID) (RegularPost, error)
}
