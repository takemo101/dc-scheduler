package domain

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"path"

	"github.com/google/uuid"
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

// --- HashPassword ValueObject ---

// HashPassword Adminのハッシュパスワード
type HashPassword []byte

// NewHashPasswordFromPlainText プレーンテキストからVOを生成するコンストラクタ
func NewHashPasswordFromPlainText(plainText string) (vo HashPassword, err error) {

	length := len(plainText)
	// 3文字以上20文字以下のパスワードしか設定できない
	if length < 3 || length > 20 {
		return vo, errors.New("HashPasswordは3-20文字の間で設定してください")
	}

	// パスワードをハッシュに変換
	hash, hashError := CreateHashPassword([]byte(plainText))

	if hashError != nil {
		return vo, hashError
	}

	return NewHashPassword(hash), err
}

// NewHashPassword コンストラクタ
func NewHashPassword(hash []byte) HashPassword {
	return HashPassword(hash)
}

// Value 値を返す
func (vo HashPassword) Value() []byte {
	return []byte(vo)
}

// Compare パスワードが一致するか
func (vo HashPassword) Compare(plainPass string) bool {
	return CompareHashPassword(vo.Value(), plainPass)
}

// --- UUID ---

// UUID VOのID VOに保持して利用する
type UUID string

// NewUUID コンストラクタ
func NewUUID(uuid string) (UUID, error) {
	if len(uuid) == 0 {
		return "", errors.New("UUIDがありません")
	}

	return UUID(uuid), nil
}

// GenerateUUID 自動生成コンストラクタ
func GenerateUUID() UUID {
	return UUID(uuid.NewString())
}

// Value IDの値を返す
func (vo UUID) Value() string {
	return string(vo)
}

// Equals VOの値が一致するか
func (vo UUID) Equals(eq UUID) bool {
	return vo.Value() == eq.Value()
}

// --- FilePath ---

// FilePath ファイルパスVO
type FilePath struct {
	directory string
	name      string
}

// NewFilePath コンストラクタ
func NewFilePath(
	directory string,
	name string,
) FilePath {
	if len(name) == 0 {
		return GenerateFilePath(directory)
	}

	return FilePath{
		directory,
		name,
	}
}

// GenerateFilePath 自動生成コンストラクタ
func GenerateFilePath(directory string) FilePath {
	name := uuid.NewString()

	return FilePath{
		directory,
		name,
	}
}

// Directory ディレクトリー
func (vo FilePath) Directory() string {
	return vo.directory
}

// Name ファイル名
func (vo FilePath) Name() string {
	return vo.name
}

// Value フルパスの値を返す
func (vo FilePath) Value() string {

	if len(vo.directory) > 0 {
		return path.Join(
			vo.directory,
			vo.name,
		)
	}

	return vo.name
}

// --- GenerateRandomKey ---

// ランダムな文字列を生成する
func GenerateRandomString(length uint) (key string, err error) {
	b := make([]byte, length)
	_, err = rand.Read(b)
	if err != nil {
		return key, err
	}

	return base64.RawURLEncoding.EncodeToString(b), err
}
