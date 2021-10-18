package domain

import (
	"errors"
	"unicode/utf8"
)

// --- AdminID ValueObject ---

// AdminID AdminのID
type AdminID struct {
	ID Identity
}

// NewAdminID コンストラクタ
func NewAdminID(id uint) (vo AdminID, err error) {

	if identity, err := NewIdentity(id); err == nil {
		return AdminID{
			ID: identity,
		}, err
	}

	return vo, err
}

// Value IDの値を返す
func (vo AdminID) Value() uint {
	return vo.ID.Value()
}

// Equals VOの値が一致するか
func (vo AdminID) Equals(eq AdminID) bool {
	return vo.Value() == eq.Value()
}

// --- AdminName ValueObject ---

// AdminName Adminの名前
type AdminName string

// NewAdminName コンストラクタ
func NewAdminName(name string) (vo AdminName, err error) {
	length := utf8.RuneCountInString(name)

	// 100文字以上は設定できない
	if length == 0 || length > 100 {
		return vo, errors.New("AdminNameは100文字以内で設定してください")
	}

	return AdminName(name), err
}

// Value 値を返す
func (vo AdminName) Value() string {
	return string(vo)
}

// --- AdminEmail ValueObject ---

// AdminEmail Adminのメールアドレス
type AdminEmail string

// NewAdminEmail コンストラクタ
func NewAdminEmail(email string) (vo AdminEmail, err error) {
	length := len(email)

	// 180文字以上は設定できない
	if length == 0 || length > 180 {
		return vo, errors.New("AdminEmailは180文字以内で設定してください")
	}

	return AdminEmail(email), err
}

// Value 値を返す
func (vo AdminEmail) Value() string {
	return string(vo)
}

// Equals VOの値が一致するか
func (vo AdminEmail) Equals(eq AdminEmail) bool {
	return vo.Value() == eq.Value()
}

// --- AdminRole ValueObject ---

// AdminRole Adminのロール
type AdminRole string

const (
	AdminRoleSystem AdminRole = "system"
	AdminRoleNormal AdminRole = "normal"
)

// NewAdminRole コンストラクタ
func NewAdminRole(role string) (vo AdminRole, err error) {

	vo = AdminRole(role)

	if !vo.Valid() {
		return vo, errors.New("AdminRoleの値が不正です")
	}

	return vo, err
}

// Value 値を返す
func (vo AdminRole) Value() string {
	return vo.String()
}

// String stringへの変換
func (vo AdminRole) String() string {
	return string(vo)
}

// Name ロールの日本語名を返す
func (vo AdminRole) Name() string {
	switch vo {
	case AdminRoleSystem:
		return "システム管理者"
	case AdminRoleNormal:
		return "通常管理者"
	}
	return ""
}

// Valid 定義したものに一致するか
func (vo AdminRole) Valid() bool {
	switch vo {
	case AdminRoleSystem, AdminRoleNormal:
		return true
	}
	return false
}

// Equals VOの値が一致するか
func (vo AdminRole) Equals(eq AdminRole) bool {
	return vo.Value() == eq.Value()
}

// AdminRoleToArray ロールをキー値形式で返す
func AdminRoleToArray() []KeyValue {
	return []KeyValue{
		{
			Key:   string(AdminRoleSystem),
			Value: AdminRoleSystem.Name(),
		},
		{
			Key:   string(AdminRoleNormal),
			Value: AdminRoleNormal.Name(),
		},
	}
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

// Equals パスワードが一致するか
func (vo HashPassword) Compare(plainPass string) bool {
	return CompareHashPassword(vo.Value(), plainPass)
}

// --- AdminAuth Entity ---

// AdminAuth 認証Entity
type AdminAuth struct {
	id    AdminID
	name  AdminName
	email AdminEmail
	role  AdminRole
}

// NewAdminAuth コンストラクタ
func NewAdminAuth(
	id uint,
	name string,
	email string,
	role AdminRole,
) AdminAuth {
	return AdminAuth{
		id: AdminID{
			ID: Identity(id),
		},
		name:  AdminName(name),
		email: AdminEmail(email),
		role:  role,
	}
}

// HaveRole 指定したロールの権限を持っているか？
func (entity AdminAuth) HaveRole(role string) bool {
	return entity.Role().Equals(AdminRole(role))
}

func (entity AdminAuth) ID() AdminID {
	return entity.id
}

func (entity AdminAuth) Name() AdminName {
	return entity.name
}

func (entity AdminAuth) Email() AdminEmail {
	return entity.email
}

func (entity AdminAuth) Role() AdminRole {
	return entity.role
}

// --- Admin Entity ---

// Admin 管理者Entity
type Admin struct {
	id           AdminID
	name         AdminName
	email        AdminEmail
	hashPassword HashPassword
	role         AdminRole
}

// NewAdmin コンストラクタ
func NewAdmin(
	id uint,
	name string,
	email string,
	hashPass []byte,
	role AdminRole,
) Admin {
	return Admin{
		id: AdminID{
			ID: Identity(id),
		},
		name:         AdminName(name),
		email:        AdminEmail(email),
		hashPassword: NewHashPassword(hashPass),
		role:         role,
	}
}

// CreateAdmin Adminを追加
func CreateAdmin(
	id uint,
	name string,
	email string,
	plainPass string,
	role string,
) (entity Admin, err error) {

	idVO, err := NewAdminID(id)
	if err != nil {
		return entity, err
	}

	nameVO, err := NewAdminName(name)
	if err != nil {
		return entity, err
	}

	emailVO, err := NewAdminEmail(email)
	if err != nil {
		return entity, err
	}

	hassVO, err := NewHashPasswordFromPlainText(plainPass)
	if err != nil {
		return entity, err
	}

	roleVO, err := NewAdminRole(role)
	if err != nil {
		return entity, err
	}

	return Admin{
		id:           idVO,
		name:         nameVO,
		email:        emailVO,
		hashPassword: hassVO,
		role:         roleVO,
	}, err
}

// Update Adminを更新
func (entity *Admin) Update(
	name string,
	email string,
	plainPass string,
	role string,
) error {
	nameVO, err := NewAdminName(name)
	if err != nil {
		return err
	}
	entity.name = nameVO

	emailVO, err := NewAdminEmail(email)
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

	roleVO, err := NewAdminRole(role)
	if err != nil {
		return err
	}
	entity.role = roleVO

	return err
}

func (entity Admin) ID() AdminID {
	return entity.id
}

func (entity Admin) Name() AdminName {
	return entity.name
}

func (entity Admin) Email() AdminEmail {
	return entity.email
}

func (entity Admin) HashPassword() HashPassword {
	return entity.hashPassword
}

func (entity Admin) Role() AdminRole {
	return entity.role
}

// CreateLoginAuth AdminからログインのためにAdminAuthを作成
func (entity Admin) CreateLoginAuth() AdminAuth {
	return AdminAuth{
		id:    entity.ID(),
		name:  entity.Name(),
		email: entity.Email(),
		role:  entity.Role(),
	}
}

// ComparePassword パスワードが一致するか
func (entity Admin) ComparePassword(plainPass string) (ok bool, err error) {
	return entity.HashPassword().Compare(plainPass), err
}

// Equals Entityが同一か
func (entity Admin) Equals(eq Admin) bool {
	return entity.ID().Equals(eq.ID()) && entity.Email().Equals(eq.Email())
}

// --- AdminService ---

// AdminService 管理者ドメインサービス
type AdminService struct {
	repository AdminRepository
}

// NewAdminService コンストラクタ
func NewAdminService(
	repository AdminRepository,
) AdminService {
	return AdminService{
		repository,
	}
}

// IsDuplicate Adminが重複しているか
func (service AdminService) IsDuplicate(admin Admin) (bool, error) {
	return service.repository.ExistsByEmail(admin.Email())
}

// IsDuplicate 指定のAdminを除き重複しているか
func (service AdminService) IsDuplicateWithoutSelf(admin Admin) (bool, error) {
	return service.repository.ExistsByIDEmail(admin.ID(), admin.Email())
}

// --- AdminRepository ---

// AdminRepository 管理者Entityの永続化
type AdminRepository interface {
	Store(entity Admin) (AdminID, error)
	Update(entity Admin) error
	FindByID(id AdminID) (Admin, error)
	FindByEmail(email AdminEmail) (Admin, error)
	Delete(id AdminID) error
	ExistsByEmail(email AdminEmail) (bool, error)
	ExistsByIDEmail(id AdminID, email AdminEmail) (bool, error)
	NextIdentity() (AdminID, error)
}

// --- AdminAuthContext ---

// AdminAuthContext コンテキスト単位での管理者認証
type AdminAuthContext interface {
	Login(entity AdminAuth) error
	Logout() error
	AdminAuth() (AdminAuth, error)
	IsLogin() bool
}
