package domain

import (
	"errors"
	"time"
)

// --- MessageApiKey ValueObject ---

// MessageApiKey ApiPostを参照するためのキー
type MessageApiKey string

// NewMessageApiKey コンストラクタ
func NewMessageApiKey(key string) (vo MessageApiKey, err error) {
	vo = MessageApiKey(key)

	if vo.IsEmpty() {
		return vo, errors.New("MessageApiKeyが空です")
	}

	return vo, err
}

// GenerateMessageApiKey MessageApiKeyの生成
func GenerateMessageApiKey() (vo MessageApiKey, err error) {

	key, err := GenerateRandomString(64)
	if err != nil {
		return vo, err
	}

	return NewMessageApiKey(key)
}

// Value 値を返す
func (vo MessageApiKey) Value() string {
	return string(vo)
}

// IsEmpty MessageApiKeyが空か
func (vo MessageApiKey) IsEmpty() bool {
	return len(vo.Value()) == 0
}

// Equals VOの値が一致するか
func (vo MessageApiKey) Equals(eq MessageApiKey) bool {
	return vo.Value() == eq.Value()
}

// --- ApiPost Entity ---

// ApiPost Api配信メッセージEntity
type ApiPost struct {
	PostMessage
	apiKey MessageApiKey
}

// NewApiPost コンストラクタ
func NewApiPost(
	id uint,
	message string,
	apiKey string,
	bot Bot,
	sentMessages []SentMessage,
) ApiPost {
	return ApiPost{
		PostMessage: NewPostMessage(
			id,
			message,
			MessageTypeApiPost,
			bot,
			sentMessages,
		),
		apiKey: MessageApiKey(apiKey),
	}
}

// CreateApiPost Api配信メッセージを生成する
func CreateApiPost(
	id uint,
	bot Bot,
) (entity ApiPost, err error) {
	if !bot.IsActive() {
		return entity, errors.New("Botが無効となっています")
	}

	postMessage, err := CreateEmptyPostMessage(
		id,
		MessageTypeApiPost,
		bot,
	)
	if err != nil {
		return entity, err
	}

	keyVO, err := GenerateMessageApiKey()
	if err != nil {
		return entity, err
	}

	return ApiPost{
		PostMessage: postMessage,
		apiKey:      keyVO,
	}, err
}

func (entity ApiPost) ApiKey() MessageApiKey {
	return entity.apiKey
}

// Send メッセージを配信する
func (entity *ApiPost) Send(message string, now time.Time) (send SentMessage, err error) {
	err = entity.ChangeMessage(message)
	if err != nil {
		return send, err
	}

	return entity.PostMessage.Send(now)
}

// CanSent Api配信可能か
func (entity ApiPost) CanSent() bool {
	return !entity.Message().IsEmpty() && entity.CanSent()
}

// Equals Entityが同一か
func (entity ApiPost) Equals(eq ApiPost) bool {
	return entity.apiKey.Equals(eq.ApiKey()) && entity.PostMessage.Equals(eq.PostMessage)
}

// --- ApiPostRepository ---

// ApiPostRepository Api配信Entityの永続化
type ApiPostRepository interface {
	PostMessageRepository
	Store(entity ApiPost) (PostMessageID, error)
	Update(entity ApiPost) error
	FindByID(id PostMessageID) (ApiPost, error)
	FindByApiKey(apiKey MessageApiKey) (ApiPost, error)
}
