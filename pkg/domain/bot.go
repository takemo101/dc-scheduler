package domain

import (
	"errors"
	"unicode/utf8"
)

// --- BotID ValueObject ---

// BotID BotのID
type BotID struct {
	ID Identity
}

// NewBotID コンストラクタ
func NewBotID(id uint) (vo BotID, err error) {

	if identity, err := NewIdentity(id); err == nil {
		return BotID{
			ID: identity,
		}, err
	}

	return vo, err
}

// Value IDの値を返す
func (vo BotID) Value() uint {
	return vo.ID.Value()
}

// Equals VOの値が一致するか
func (vo BotID) Equals(eq BotID) bool {
	return vo.Value() == eq.Value()
}

// --- BotLabel ValueObject ---

// BotName Botの名前
type BotName string

// NewBotName コンストラクタ
func NewBotName(name string) (vo BotName, err error) {
	length := utf8.RuneCountInString(name)

	// 32文字以上は設定できない
	if length == 0 || length > 32 {
		return vo, errors.New("BotNameは32文字以内で設定してください")
	}

	return BotName(name), err
}

// Value 値を返す
func (vo BotName) Value() string {
	return string(vo)
}
