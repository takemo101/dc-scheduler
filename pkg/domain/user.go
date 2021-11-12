package domain

import (
	"errors"
	"unicode/utf8"
)

// --- UserID ValueObject ---

// UserID UserのID
type UserID struct {
	ID Identity
}

// NewUserID コンストラクタ
func NewUserID(id uint) (vo UserID, err error) {

	if identity, err := NewIdentity(id); err == nil {
		return UserID{
			ID: identity,
		}, err
	}

	return vo, err
}

// Value IDの値を返す
func (vo UserID) Value() uint {
	return vo.ID.Value()
}

// Equals VOの値が一致するか
func (vo UserID) Equals(eq UserID) bool {
	return vo.Value() == eq.Value()
}

// --- UserName ValueObject ---

// UserName Userの名前
type UserName string

// NewUserName コンストラクタ
func NewUserName(name string) (vo UserName, err error) {
	length := utf8.RuneCountInString(name)

	// 100文字以上は設定できない
	if length == 0 || length > 100 {
		return vo, errors.New("UserNameは100文字以内で設定してください")
	}

	return UserName(name), err
}

// Value 値を返す
func (vo UserName) Value() string {
	return string(vo)
}

// --- UserEmail ValueObject ---

// UserEmail Userのメールアドレス
type UserEmail string

// NewUserEmail コンストラクタ
func NewUserEmail(email string) (vo UserEmail, err error) {
	length := len(email)

	// 180文字以上は設定できない
	if length == 0 || length > 180 {
		return vo, errors.New("UserEmailは180文字以内で設定してください")
	}

	return UserEmail(email), err
}

// Value 値を返す
func (vo UserEmail) Value() string {
	return string(vo)
}

// Equals VOの値が一致するか
func (vo UserEmail) Equals(eq UserEmail) bool {
	return vo.Value() == eq.Value()
}

// --- UserActivationKey ValueObject ---

// UserActivationKey UserのActivationKey
type UserActivationKey string

// GenerateUserActivationKey ActivationKeyの生成
func GenerateUserActivationKey() (vo UserActivationKey, err error) {

	key, err := GenerateRandomString(32)
	if err != nil {
		return vo, err
	}

	vo = UserActivationKey(key)

	if vo.IsEmpty() {
		return vo, errors.New("UserActivationKeyの生成に失敗しました")
	}

	return vo, err
}

// Value 値を返す
func (vo UserActivationKey) Value() string {
	return string(vo)
}

// IsEmpty ActivationKeyが空か
func (vo UserActivationKey) IsEmpty() bool {
	return len(vo.Value()) == 0
}

// Equals VOの値が一致するか
func (vo UserActivationKey) Equals(eq UserActivationKey) bool {
	return vo.Value() == eq.Value()
}

// --- UserAuth Entity ---

// UserAuth 認証Entity
type UserAuth struct {
	id    UserID
	name  UserName
	email UserEmail
}

// NewUserAuth コンストラクタ
func NewUserAuth(
	id uint,
	name string,
	email string,
) UserAuth {
	return UserAuth{
		id: UserID{
			ID: Identity(id),
		},
		name:  UserName(name),
		email: UserEmail(email),
	}
}

func (entity UserAuth) ID() UserID {
	return entity.id
}

func (entity UserAuth) Name() UserName {
	return entity.name
}

func (entity UserAuth) Email() UserEmail {
	return entity.email
}

// --- UserSignature DTO ---

// UserSignature シグネチャDTO
type UserSignature struct {
	Email         UserEmail         `json:"email"`
	ActivationKey UserActivationKey `json:"activation_key"`
}

// --- User Entity ---

// User 管理者Entity
type User struct {
	id            UserID
	name          UserName
	email         UserEmail
	hashPassword  HashPassword
	activationKey UserActivationKey
	active        bool
}

// NewUser コンストラクタ
func NewUser(
	id uint,
	name string,
	email string,
	hashPass []byte,
	activationKey string,
	active bool,
) User {
	return User{
		id: UserID{
			ID: Identity(id),
		},
		name:          UserName(name),
		email:         UserEmail(email),
		hashPassword:  NewHashPassword(hashPass),
		activationKey: UserActivationKey(activationKey),
		active:        active,
	}
}

// CreatedUser Userを追加（メール認証無視）
func CreateUser(
	id uint,
	name string,
	email string,
	plainPass string,
	active bool,
) (entity User, err error) {

	idVO, err := NewUserID(id)
	if err != nil {
		return entity, err
	}

	nameVO, err := NewUserName(name)
	if err != nil {
		return entity, err
	}

	emailVO, err := NewUserEmail(email)
	if err != nil {
		return entity, err
	}

	hassVO, err := NewHashPasswordFromPlainText(plainPass)
	if err != nil {
		return entity, err
	}

	keyVO, err := GenerateUserActivationKey()
	if err != nil {
		return entity, err
	}

	return User{
		id:            idVO,
		name:          nameVO,
		email:         emailVO,
		hashPassword:  hassVO,
		activationKey: keyVO,
		active:        active,
	}, err
}

// CreateTemporaryUser Userを仮追加
func CreateTemporaryUser(
	id uint,
	name string,
	email string,
	plainPass string,
) (entity User, err error) {

	idVO, err := NewUserID(id)
	if err != nil {
		return entity, err
	}

	nameVO, err := NewUserName(name)
	if err != nil {
		return entity, err
	}

	emailVO, err := NewUserEmail(email)
	if err != nil {
		return entity, err
	}

	hassVO, err := NewHashPasswordFromPlainText(plainPass)
	if err != nil {
		return entity, err
	}

	keyVO, err := GenerateUserActivationKey()
	if err != nil {
		return entity, err
	}

	return User{
		id:            idVO,
		name:          nameVO,
		email:         emailVO,
		hashPassword:  hassVO,
		activationKey: keyVO,
		active:        false,
	}, err
}

// Update Userを更新
func (entity *User) Update(
	name string,
	email string,
	plainPass string,
	active bool,
) error {
	nameVO, err := NewUserName(name)
	if err != nil {
		return err
	}
	entity.name = nameVO

	emailVO, err := NewUserEmail(email)
	if err != nil {
		return err
	}
	entity.email = emailVO

	if len(plainPass) > 0 {
		hassVO, err := NewHashPasswordFromPlainText(plainPass)
		if err != nil {
			return err
		}
		entity.hashPassword = hassVO
	}

	entity.active = active

	return err
}

func (entity User) ID() UserID {
	return entity.id
}

func (entity User) Name() UserName {
	return entity.name
}

func (entity User) Email() UserEmail {
	return entity.email
}

func (entity User) HashPassword() HashPassword {
	return entity.hashPassword
}

func (entity User) ActivationKey() UserActivationKey {
	return entity.activationKey
}

func (entity User) IsActivated() bool {
	return entity.active
}

// Activation アクティベートを行う
func (entity *User) Activation(activationKey UserActivationKey) (ok bool, err error) {

	if entity.IsActivated() {
		return false, errors.New("Activation済みです")
	}

	if entity.ActivationKey().IsEmpty() {
		return false, errors.New("ActivationKeyが空です")
	}

	// キーの一致をチェックする
	if entity.ActivationKey().Equals(activationKey) {
		return false, err
	}

	entity.active = true

	return true, err
}

// CreateLoginAuth UserからログインのためにUserAuthを作成
func (entity User) CreateLoginAuth() (auth UserAuth, err error) {

	if !entity.IsActivated() {
		return auth, errors.New("Activationが済んでいません")
	}

	return UserAuth{
		id:    entity.ID(),
		name:  entity.Name(),
		email: entity.Email(),
	}, err
}

// CreateSignature Userからシグネチャ生成のためにUserSignatureを作成
func (entity User) CreateSignature() UserSignature {
	return UserSignature{
		Email:         entity.Email(),
		ActivationKey: entity.ActivationKey(),
	}
}

// ComparePassword パスワードが一致するか
func (entity User) ComparePassword(plainPass string) (ok bool, err error) {
	return entity.HashPassword().Compare(plainPass), err
}

// Equals Entityが同一か
func (entity User) Equals(eq User) bool {
	return entity.ID().Equals(eq.ID()) && entity.Email().Equals(eq.Email())
}

// --- UserService ---

// UserService 管理者ドメインサービス
type UserService struct {
	repository UserRepository
	crypter    UserActivationSignatureCrypter
}

// NewUserService コンストラクタ
func NewUserService(
	repository UserRepository,
	crypter UserActivationSignatureCrypter,
) UserService {
	return UserService{
		repository,
		crypter,
	}
}

// IsDuplicate Userが重複しているか
func (service UserService) IsDuplicate(admin User) (bool, error) {
	return service.repository.ExistsByEmail(admin.Email())
}

// IsDuplicate 指定のUserを除き重複しているか
func (service UserService) IsDuplicateWithoutSelf(admin User) (bool, error) {
	return service.repository.ExistsByIDEmail(admin.ID(), admin.Email())
}

// --- UserActivationNotify ---

// UserActivationNotify User仮登録時のアクティベーションメール通知
type UserActivationNotify struct {
	user      User
	signature string
}

// CreateUserActivationNotify アクティベーションメール通知の生成
func CreateUserActivationNotify(
	user User,
	signature string,
) (notify TemplateMailNotify, err error) {

	if user.IsActivated() {
		return notify, errors.New("UserがActivation済みにて通知できません")
	}

	return UserActivationNotify{
		user,
		signature,
	}, err
}

// ToEmailAddress 送信先アドレス
func (notify UserActivationNotify) ToEmailAddress() string {
	return notify.user.Email().Value()
}

// Key 識別子
func (notify UserActivationNotify) Key() string {
	return "user-activation"
}

// Data テンプレートに与えるデータ
func (notify UserActivationNotify) Data() TemplateData {
	return TemplateData{
		"name":      notify.user.Name().Value(),
		"email":     notify.user.Email().Value(),
		"signature": notify.signature,
	}
}

// --- UserRepository ---

// UserRepository 管理者Entityの永続化
type UserRepository interface {
	Store(entity User) (UserID, error)
	Update(entity User) error
	FindByID(id UserID) (User, error)
	FindByEmail(email UserEmail) (User, error)
	Delete(id UserID) error
	ExistsByEmail(email UserEmail) (bool, error)
	ExistsByIDEmail(id UserID, email UserEmail) (bool, error)
	NextIdentity() (UserID, error)
}

// --- UserAuthContext ---

// UserAuthContext コンテキスト単位での管理者認証
type UserAuthContext interface {
	Login(entity UserAuth) error
	Logout() error
	UserAuth() (UserAuth, error)
	IsLogin() bool
}

// --- UserActivationSignatureCrypter ---

// UserActivationSignatureCrypter アクティベーションのためのシグネチャ暗号複合化
type UserActivationSignatureCrypter interface {
	Encrypt(dto UserSignature) (string, error)
	Decrypt(signature string) (UserSignature, error)
}
