package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// --- Identity ---

// Identity VOのID VOに保持して利用する
type Identity uint

// NewIdentity コンストラクタ
func NewIdentity(id uint) (vo Identity, err error) {
	if id == 0 {
		return vo, errors.New("Identityが不正です")
	}
	return Identity(id), err
}

// Value IDの値を返す
func (vo Identity) Value() uint {
	return uint(vo)
}

// --- KeyValue ---

// キーと値を表現した構造体
type KeyValue struct {
	Key   interface{} `json:"key"`
	Value string      `json:"value"`
}

// --- HashPassword ---

// CreateHashPassword ハッシュパスワードを生成
func CreateHashPassword(plainPass []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(plainPass, bcrypt.DefaultCost)
	if err != nil {
		return hash, err
	}

	return hash, nil
}

// CompareHashPassword ハッシュパスワードを比較
func CompareHashPassword(hash []byte, plainPass string) bool {
	bytePass := []byte(plainPass)

	// check
	if err := bcrypt.CompareHashAndPassword(hash, bytePass); err != nil {
		return false
	}

	return true
}
