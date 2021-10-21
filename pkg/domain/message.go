package domain

import (
	"errors"
	"time"
	"unicode/utf8"
)

// --- PostMessageID ValueObject ---

// PostMessageID PostMessageのID
type PostMessageID struct {
	ID Identity
}

// NewPostMessageID コンストラクタ
func NewPostMessageID(id uint) (vo PostMessageID, err error) {

	if identity, err := NewIdentity(id); err == nil {
		return PostMessageID{
			ID: identity,
		}, err
	}

	return vo, err
}

// Value IDの値を返す
func (vo PostMessageID) Value() uint {
	return vo.ID.Value()
}

// Equals VOの値が一致するか
func (vo PostMessageID) Equals(eq PostMessageID) bool {
	return vo.Value() == eq.Value()
}

// --- Message ValueObject ---

// Message メッセージ
type Message string

// NewMessage コンストラクタ
func NewMessage(message string) (vo Message, err error) {
	length := utf8.RuneCountInString(message)

	// 2000文字以上は設定できない
	if length == 0 || length > 2000 {
		return vo, errors.New("Messageは2000文字以内で設定してください")
	}

	return Message(message), err
}

// Value 値を返す
func (vo Message) Value() string {
	return string(vo)
}

// --- MessageType ValueObject ---

// MessageType メッセージタイプ
type MessageType string

const (
	MessageTypeSchedulePost  MessageType = "schedule"
	MessageTypeImmediatePost MessageType = "immediate"
	MessageTypeRegularPost   MessageType = "regular"
)

// NewMessageType コンストラクタ
func NewMessageType(messageType string) (vo MessageType, err error) {

	vo = MessageType(messageType)

	if !vo.Valid() {
		return vo, errors.New("MessageTypeの値が不正です")
	}

	return vo, err
}

// Value 値を返す
func (vo MessageType) Value() string {
	return vo.String()
}

// String stringへの変換
func (vo MessageType) String() string {
	return string(vo)
}

// Name メッセージタイプの日本語名を返す
func (vo MessageType) Name() string {
	switch vo {
	case MessageTypeSchedulePost:
		return "予約配信"
	case MessageTypeImmediatePost:
		return "即時配信"
	case MessageTypeRegularPost:
		return "定期配信"
	}
	return ""
}

// Valid 定義したものに一致するか
func (vo MessageType) Valid() bool {
	switch vo {
	case MessageTypeSchedulePost, MessageTypeImmediatePost, MessageTypeRegularPost:
		return true
	}
	return false
}

// Equals VOの値が一致するか
func (vo MessageType) Equals(eq MessageType) bool {
	return vo.Value() == eq.Value()
}

// --- MessageSendedAt ValueObject ---

// MessageSendedAt 送信日時
type MessageSendedAt time.Time

// NewMessageSendedAt コンストラクタ
func NewMessageSendedAt(at time.Time) MessageSendedAt {
	return MessageSendedAt(at)
}

// Value 値を返す
func (vo MessageSendedAt) Value() time.Time {
	return time.Time(vo)
}

// --- SentMessage Entity ---

// SentMessage メッセージ送信（履歴）Entity
type SentMessage struct {
	id       UUID
	message  Message
	sendedAt MessageSendedAt
}

// SendMessage メッセージ送信をする
func SendMessage(
	messageVO Message,
	now time.Time,
) (entity SentMessage, err error) {
	sendedAtVO := NewMessageSendedAt(now)

	return SentMessage{
		id:       GenerateUUID(),
		message:  messageVO,
		sendedAt: sendedAtVO,
	}, err
}

func (entity SentMessage) Message() Message {
	return entity.message
}

func (entity SentMessage) SendedAt() MessageSendedAt {
	return entity.sendedAt
}

// --- PostMessage Entity ---

// PostMessage 配信メッセージEntity
type PostMessage struct {
	id           PostMessageID
	message      Message
	messageType  MessageType
	bot          Bot
	sentMessages []SentMessage
}

// NewPostMessage コンストラクタ
func NewPostMessage(
	id uint,
	message string,
	messageType MessageType,
	bot Bot,
	sentMessages []SentMessage,
) PostMessage {
	return PostMessage{
		id: PostMessageID{
			ID: Identity(id),
		},
		message:      Message(message),
		messageType:  messageType,
		bot:          bot,
		sentMessages: sentMessages,
	}
}

// CreatePostMessage 配信メッセージを生成する
func CreatePostMessage(
	id uint,
	message string,
	messageType MessageType,
	bot Bot,
) (entity PostMessage, err error) {
	idVO, err := NewPostMessageID(id)
	if err != nil {
		return entity, err
	}

	messageVO, err := NewMessage(message)
	if err != nil {
		return entity, err
	}

	messageTypeVO, err := NewMessageType(messageType.Value())
	if err != nil {
		return entity, err
	}

	return PostMessage{
		id:           idVO,
		message:      messageVO,
		messageType:  messageTypeVO,
		bot:          bot,
		sentMessages: []SentMessage{},
	}, err
}

func (entity PostMessage) ID() PostMessageID {
	return entity.id
}

func (entity PostMessage) Message() Message {
	return entity.message
}

func (entity PostMessage) MessageType() MessageType {
	return entity.messageType
}

func (entity PostMessage) Bot() Bot {
	return entity.bot
}

func (entity PostMessage) SentMessages() []SentMessage {
	return entity.sentMessages
}

// Send メッセージを送信する
func (entity *PostMessage) Send(now time.Time) (send SentMessage, err error) {
	if !entity.CanSent() {
		return send, errors.New("Message送信可能ではありません")
	}

	send, err = SendMessage(entity.message, now)

	entity.sentMessages = append(
		entity.sentMessages,
		send,
	)

	return send, err
}

// ChangeMessage メッセージを変更する
func (entity PostMessage) ChangeMessage(message string) (err error) {
	messageVO, err := NewMessage(message)
	if err != nil {
		return err
	}
	entity.message = messageVO

	return err
}

// CanSent メッセージが送信可能か
func (entity PostMessage) CanSent() bool {
	return entity.bot.IsActive()
}

// IsSended メッセージ送信済みの状態か
func (entity PostMessage) IsSended() bool {
	return len(entity.sentMessages) > 0
}

// Equals Entityが同一か
func (entity PostMessage) Equals(eq PostMessage) bool {
	return entity.ID().Equals(
		eq.ID(),
	) && entity.MessageType().Equals(
		eq.MessageType(),
	) && entity.bot.Equals(
		eq.Bot(),
	)
}

// --- ImmediatePost Entity ---

// ImmediatePost 即時配信メッセージEntity
type ImmediatePost struct {
	PostMessage
}

// NewImmediatePost コンストラクタ
func NewImmediatePost(
	id uint,
	message string,
	bot Bot,
	sentMessages []SentMessage,
) ImmediatePost {
	return ImmediatePost{
		PostMessage: NewPostMessage(
			id,
			message,
			MessageTypeImmediatePost,
			bot,
			sentMessages,
		),
	}
}

// CreateImmediatePost 即時配信メッセージを生成する
func CreateImmediatePost(
	id uint,
	message string,
	bot Bot,
) (entity ImmediatePost, err error) {
	postMessage, err := CreatePostMessage(
		id,
		message,
		MessageTypeImmediatePost,
		bot,
	)

	return ImmediatePost{
		PostMessage: postMessage,
	}, err
}

// --- ReservationAt ValueObject ---

// MessageReservationAt 配信予約日時
type MessageReservationAt time.Time

// NewMessageReservationAt コンストラクタ
func NewMessageReservationAt(at time.Time, now time.Time) (vo MessageReservationAt, err error) {
	if at.After(now) {
		return MessageReservationAt(at), err
	}

	return vo, errors.New("MessageReservationAtは現在以降を指定してください")
}

// Value 値を返す
func (vo MessageReservationAt) Value() time.Time {
	return time.Time(vo)
}

// After 予約日時以降か
func (vo MessageReservationAt) After(now time.Time) bool {
	return !vo.Value().Before(now)
}

// --- SchedulePost Entity ---

// SchedulePost 予約配信メッセージEntity
type SchedulePost struct {
	PostMessage
	reservationAt MessageReservationAt
}

// NewSchedulePost コンストラクタ
func NewSchedulePost(
	id uint,
	message string,
	reservationAt time.Time,
	bot Bot,
	sentMessages []SentMessage,
) SchedulePost {
	return SchedulePost{
		PostMessage: NewPostMessage(
			id,
			message,
			MessageTypeSchedulePost,
			bot,
			sentMessages,
		),
		reservationAt: MessageReservationAt(reservationAt),
	}
}

// CreateSchedulePost 即時配信メッセージを生成する
func CreateSchedulePost(
	id uint,
	message string,
	reservationAt time.Time,
	bot Bot,
	now time.Time,
) (entity SchedulePost, err error) {
	postMessage, err := CreatePostMessage(
		id,
		message,
		MessageTypeSchedulePost,
		bot,
	)

	reservationAtVO, err := NewMessageReservationAt(reservationAt, now)

	return SchedulePost{
		PostMessage:   postMessage,
		reservationAt: reservationAtVO,
	}, err
}

func (entity SchedulePost) ReservationAt() MessageReservationAt {
	return entity.reservationAt
}

// IsSended 予約時間を過ぎているか
func (entity SchedulePost) IsPassedReservationAt(now time.Time) bool {
	return entity.ReservationAt().After(now)
}

// CanSent 予約配信可能か
func (entity SchedulePost) CanSent(now time.Time) bool {
	return entity.Bot().IsActive() && entity.IsPassedReservationAt(now)
}

// Send メッセージを送信する
func (entity *SchedulePost) Send(now time.Time) (send SentMessage, err error) {
	if !entity.CanSent(now) {
		return send, errors.New("Message送信可能ではありません")
	}

	send, err = SendMessage(entity.message, now)

	entity.sentMessages = append(
		entity.sentMessages,
		send,
	)

	return send, err
}

// --- PostMessageRepository ---

type PostMessageRepository interface {
	SaveSendedMessage(id PostMessageID, entity SentMessage) error
	Delete(id PostMessageID) error
	NextIdentity() (PostMessageID, error)
}

// --- ImmediatePostRepository ---

// ImmediatePostRepository 即時配信Entityの永続化
type ImmediatePostRepository interface {
	PostMessageRepository
	Store(entity ImmediatePost) (PostMessageID, error)
	FindByID(id PostMessageID) (ImmediatePost, error)
}

// --- SchedulePostRepository ---

// SchedulePostRepository 予約配信Entityの永続化
type SchedulePostRepository interface {
	PostMessageRepository
	SendList(now time.Time) []SchedulePost
	Store(entity SchedulePost) (PostMessageID, error)
	Update(entity SchedulePost) error
	FindByID(id PostMessageID) (SchedulePost, error)
}

// --- DiscordMessageAdapter ---

// DiscordMessageAdapter Discordメッセージを配信するアダプタ
type DiscordMessageAdapter interface {
	SendMessage(bot Bot, message Message) error
}
