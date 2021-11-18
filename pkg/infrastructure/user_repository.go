package infrastructure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/domain"
	"gorm.io/gorm"
)

// --- UserAuthContext Implement ---

const (
	UserSessionIDKey    string = "user.id"
	UserSessionNameKey  string = "user.name"
	UserSessionEmailKey string = "user.email"
)

// UserAuthContext domain.UserAuthContextの実装
type UserAuthContext struct {
	session *session.Session
}

// NewUserAuthContext コンストラクタ
func NewUserAuthContext(session *session.Session) domain.UserAuthContext {
	return &UserAuthContext{
		session: session,
	}
}

// Login ログイン状態にする
func (context *UserAuthContext) Login(entity domain.UserAuth) (err error) {
	context.session.Set(UserSessionIDKey, entity.ID().Value())
	context.session.Set(UserSessionNameKey, entity.Name().Value())
	context.session.Set(UserSessionEmailKey, entity.Email().Value())
	err = context.session.Save()

	if err != nil {
		return err
	}

	return err
}

// Logout ログアウト状態にする
func (context *UserAuthContext) Logout() (err error) {
	context.session.Set(UserSessionIDKey, nil)
	err = context.session.Save()

	return err
}

// UserAuth ログインしている管理者を取得する
func (context *UserAuthContext) UserAuth() (entity domain.UserAuth, err error) {
	id, ok := context.getPrimitiveID()
	if !ok {
		return entity, errors.New("IDを取得できません")
	}

	name, ok := context.getPrimitiveName()
	if !ok {
		return entity, errors.New("Nameを取得できません")
	}

	email, ok := context.getPrimitiveEmail()
	if !ok {
		return entity, errors.New("Emailを取得できません")
	}

	return domain.NewUserAuth(
		id,
		name,
		email,
	), err
}

// IsLogin ログインしているかチェックする
func (context *UserAuthContext) IsLogin() bool {
	_, ok := context.getPrimitiveID()
	return ok
}

// getPrimitiveID Sessionから保存したID値を取り出す
func (context *UserAuthContext) getPrimitiveID() (id uint, ok bool) {
	id, ok = context.session.Get(UserSessionIDKey).(uint)
	return id, ok
}

// getPrimitiveName Sessionから保存したName値を取り出す
func (context *UserAuthContext) getPrimitiveName() (name string, ok bool) {
	name, ok = context.session.Get(UserSessionNameKey).(string)
	return name, ok
}

// getPrimitiveEmail Sessionから保存したEmail値を取り出す
func (context *UserAuthContext) getPrimitiveEmail() (email string, ok bool) {
	email, ok = context.session.Get(UserSessionEmailKey).(string)
	return email, ok
}

// --- UserAuthContext Implement InMemory ---

// InMemoryUserAuthContext domain.UserAuthContextの実装：テスト
type InMemoryUserAuthContext struct {
	entity domain.UserAuth
	login  bool
}

// NewInMemoryUserAuthContext コンストラクタ
func NewInMemoryUserAuthContext() domain.UserAuthContext {
	return &InMemoryUserAuthContext{
		login: false,
	}
}

// Login ログイン状態にする
func (context *InMemoryUserAuthContext) Login(entity domain.UserAuth) (err error) {
	context.entity = entity
	context.login = true
	return err
}

// Logout ログアウト状態にする
func (context *InMemoryUserAuthContext) Logout() (err error) {
	context.login = false
	return err
}

// UserAuth ログインしている管理者を取得する
func (context *InMemoryUserAuthContext) UserAuth() (entity domain.UserAuth, err error) {
	if !context.login {
		return entity, errors.New("ログインしていません")
	}

	return context.entity, err
}

// IsLogin ログインしているかチェックする
func (context *InMemoryUserAuthContext) IsLogin() bool {
	return context.login
}

// --- UserRepository Implement ---

// NewUserRepository domain.UserRepositoryの実装
type UserRepository struct {
	db     core.Database
	config core.Config
}

// NewUserRepository コンストラクタ
func NewUserRepository(
	db core.Database,
	config core.Config,
) domain.UserRepository {
	return UserRepository{
		db,
		config,
	}
}

// Store Userの追加
func (repo UserRepository) Store(entity domain.User) (vo domain.UserID, err error) {
	model := User{}

	model.Name = entity.Name().Value()
	model.Email = entity.Email().Value()
	model.Password = entity.HashPassword().Value()
	model.ActivationKey = entity.ActivationKey().Value()
	model.Active = sql.NullBool{Bool: entity.IsActivated(), Valid: true}

	if err = repo.db.GormDB.Create(&model).Error; err != nil {
		return vo, err
	}

	return domain.NewUserID(model.ID)
}

// Update Userの更新
func (repo UserRepository) Update(entity domain.User) error {
	model := User{}

	model.ID = entity.ID().Value()
	model.Name = entity.Name().Value()
	model.Email = entity.Email().Value()
	model.Password = entity.HashPassword().Value()
	model.Active = sql.NullBool{Bool: entity.IsActivated(), Valid: true}

	return repo.db.GormDB.Updates(&model).Error
}

// FindByID UserIDからUserを取得する
func (repo UserRepository) FindByID(id domain.UserID) (entity domain.User, err error) {
	model := User{}

	if err = repo.db.GormDB.Where("id = ?", id.Value()).First(&model).Error; err != nil {
		return entity, err
	}

	return CreateUserEntityFromModel(model), err
}

// FindByEmail UserEmailからUserを取得する
func (repo UserRepository) FindByEmail(email domain.UserEmail) (entity domain.User, err error) {
	model := User{}

	if err = repo.db.GormDB.Where("email = ?", email.Value()).First(&model).Error; err != nil {
		return entity, err
	}

	return CreateUserEntityFromModel(model), err
}

// Delete UserIDからUserを削除する
func (repo UserRepository) Delete(id domain.UserID) error {
	return repo.db.GormDB.Where("id = ?", id.Value()).Delete(&User{}).Error
}

// ExistsByEmail UserEmailが重複しているUserがあるか
func (repo UserRepository) ExistsByEmail(email domain.UserEmail) (bool, error) {
	count := int64(0)
	err := repo.db.GormDB.Model(&User{}).
		Where("email = ?", email.Value()).
		Count(&count).
		Error

	return (count > 0), err
}

// ExistsByIDEmail 指定したUserIDを除きUserEmailが重複しているUserがあるか
func (repo UserRepository) ExistsByIDEmail(id domain.UserID, email domain.UserEmail) (bool, error) {
	count := int64(0)
	err := repo.db.GormDB.Model(&User{}).
		Where("id <> ? AND email = ?", id.Value(), email.Value()).
		Count(&count).
		Error

	return (count > 0), err
}

// NextIdentity 次のIDを取得する
func (repo UserRepository) NextIdentity() (domain.UserID, error) {
	var max uint

	sql := GetNextIdentitySelectSQL(repo.config.DB.Type)
	repo.db.GormDB.Model(&User{}).Select(sql).Scan(&max)

	return domain.NewUserID(max + 1)
}

// CreateUserEntityFromModel UserからEntityを生成する
func CreateUserEntityFromModel(model User) domain.User {
	return domain.NewUser(
		model.ID,
		model.Name,
		model.Email,
		model.Password,
		model.ActivationKey,
		model.Active.Bool,
	)
}

// --- UserActivationSignatureCrypter Implement ---

// UserActivationSignatureCrypter domain.UserActivationSignatureCrypterの実装
type UserActivationSignatureCrypter struct {
	config core.Config
}

// NewUserActivationSignatureCrypter コンストラクタ
func NewUserActivationSignatureCrypter(
	config core.Config,
) domain.UserActivationSignatureCrypter {
	return UserActivationSignatureCrypter{
		config,
	}
}

// Encrypt シグネチャへ暗号化
func (crypter UserActivationSignatureCrypter) Encrypt(dto domain.UserSignature) (signature string, err error) {

	block, err := aes.NewCipher(crypter.SecretKey())
	if err != nil {
		return signature, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return signature, err
	}

	// 初期化ベクトル作成
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return signature, err
	}

	// jsonに変換
	json, err := json.Marshal(dto)
	if err != nil {
		return signature, err
	}

	cipher := gcm.Seal(nil, nonce, json, nil) // 暗号化
	cipher = append(nonce, cipher...)         // 先頭に付与

	// Base64 Encode
	return base64.StdEncoding.EncodeToString(cipher), err
}

// Decrypt シグネチャを複合化
func (crypter UserActivationSignatureCrypter) Decrypt(base64signature string) (dto domain.UserSignature, err error) {
	// Base64 Decode
	cipherText, err := base64.StdEncoding.DecodeString(base64signature)
	if err != nil {
		return dto, err
	}

	block, err := aes.NewCipher(crypter.SecretKey())
	if err != nil {
		return dto, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return dto, err
	}

	// 初期化ベクトルを取り出す
	nonce := cipherText[:gcm.NonceSize()]
	signature, err := gcm.Open(nil, nonce, cipherText[gcm.NonceSize():], nil)
	if err != nil {
		return dto, err
	}

	err = json.Unmarshal(signature, &dto)
	return dto, err
}

// SecretKey シグネチャを生成するためのシークレットキーを返す
func (crypter UserActivationSignatureCrypter) SecretKey() []byte {
	key := []byte(crypter.config.App.Secret)
	length := len(key)

	// keyの長さは16 24 32 byte である必要があるため丸める
	switch {
	case length >= 32:
		key = key[0:32]
	case length < 32 && length >= 24:
		key = key[0:24]
	case length < 24 && length >= 16:
		key = key[0:16]
	default:
		// 足りないbyteを埋めるスライス
		baseKey := []byte("xxxxxxxxxxxxxxxx")
		// 足りない分を結合
		key = append(key, baseKey[length:16]...)
	}

	return key
}

// --- User ---

// User Gormモデル
type User struct {
	gorm.Model
	Name          string `gorm:"type:varchar(191);index;not null"`
	Email         string `gorm:"type:varchar(191);uniqueIndex;not null"`
	Password      []byte
	ActivationKey string       `gorm:"type:varchar(191);index;not null;default:user"`
	Active        sql.NullBool `gorm:"type:boolean;index"`
}
