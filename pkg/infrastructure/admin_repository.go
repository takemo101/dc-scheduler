package infrastructure

import (
	"errors"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/domain"
	"gorm.io/gorm"
)

// --- AdminAuthContext Implement ---

const (
	AdminSessionIDKey    string = "admin.id"
	AdminSessionNameKey  string = "admin.name"
	AdminSessionEmailKey string = "admin.email"
	AdminSessionRoleKey  string = "admin.role"
)

// AdminAuthContext domain.AdminAuthContextの実装
type AdminAuthContext struct {
	session *session.Session
}

// NewAdminAuthContext コンストラクタ
func NewAdminAuthContext(session *session.Session) domain.AdminAuthContext {
	return &AdminAuthContext{
		session: session,
	}
}

// Login ログイン状態にする
func (context *AdminAuthContext) Login(entity domain.AdminAuth) (err error) {
	context.session.Set(AdminSessionIDKey, entity.ID().Value())
	context.session.Set(AdminSessionNameKey, entity.Name().Value())
	context.session.Set(AdminSessionEmailKey, entity.Email().Value())
	context.session.Set(AdminSessionRoleKey, entity.Role().Value())
	err = context.session.Save()

	if err != nil {
		return err
	}

	return err
}

// Logout ログアウト状態にする
func (context *AdminAuthContext) Logout() (err error) {
	context.session.Set(AdminSessionIDKey, nil)
	context.session.Set(AdminSessionRoleKey, nil)
	err = context.session.Save()

	return err
}

// AdminAuth ログインしている管理者を取得する
func (context *AdminAuthContext) AdminAuth() (entity domain.AdminAuth, err error) {
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

	role, ok := context.getPrimitiveRole()
	if !ok {
		return entity, errors.New("Roleを取得できません")
	}

	return domain.NewAdminAuth(
		id,
		name,
		email,
		domain.AdminRole(role),
	), err
}

// IsLogin ログインしているかチェックする
func (context *AdminAuthContext) IsLogin() bool {
	_, ok := context.getPrimitiveID()
	return ok
}

// getPrimitiveID Sessionから保存したID値を取り出す
func (context *AdminAuthContext) getPrimitiveID() (id uint, ok bool) {
	id, ok = context.session.Get(AdminSessionIDKey).(uint)
	return id, ok
}

// getPrimitiveName Sessionから保存したName値を取り出す
func (context *AdminAuthContext) getPrimitiveName() (name string, ok bool) {
	name, ok = context.session.Get(AdminSessionNameKey).(string)
	return name, ok
}

// getPrimitiveEmail Sessionから保存したEmail値を取り出す
func (context *AdminAuthContext) getPrimitiveEmail() (email string, ok bool) {
	email, ok = context.session.Get(AdminSessionEmailKey).(string)
	return email, ok
}

// getPrimitiveRole Sessionから保存したRole値を取り出す
func (context *AdminAuthContext) getPrimitiveRole() (role string, ok bool) {
	role, ok = context.session.Get(AdminSessionRoleKey).(string)
	return role, ok
}

// --- AdminAuthContext Implement InMemory ---

// InMemoryAdminAuthContext domain.AdminAuthContextの実装：テスト
type InMemoryAdminAuthContext struct {
	entity domain.AdminAuth
	login  bool
}

// NewInMemoryAdminAuthContext コンストラクタ
func NewInMemoryAdminAuthContext() domain.AdminAuthContext {
	return &InMemoryAdminAuthContext{
		login: false,
	}
}

// Login ログイン状態にする
func (context *InMemoryAdminAuthContext) Login(entity domain.AdminAuth) (err error) {
	context.entity = entity
	context.login = true
	return err
}

// Logout ログアウト状態にする
func (context *InMemoryAdminAuthContext) Logout() (err error) {
	context.login = false
	return err
}

// AdminAuth ログインしている管理者を取得する
func (context *InMemoryAdminAuthContext) AdminAuth() (entity domain.AdminAuth, err error) {
	if !context.login {
		return entity, errors.New("ログインしていません")
	}

	return context.entity, err
}

// IsLogin ログインしているかチェックする
func (context *InMemoryAdminAuthContext) IsLogin() bool {
	return context.login
}

// --- AdminRepository Implement ---

// NewAdminRepository domain.AdminRepositoryの実装
type AdminRepository struct {
	db     core.Database
	config core.Config
}

// NewAdminRepository コンストラクタ
func NewAdminRepository(
	db core.Database,
	config core.Config,
) domain.AdminRepository {
	return AdminRepository{
		db,
		config,
	}
}

// Store Adminの追加
func (repo AdminRepository) Store(entity domain.Admin) (vo domain.AdminID, err error) {
	model := Admin{}

	model.Name = entity.Name().Value()
	model.Email = entity.Email().Value()
	model.Password = entity.HashPassword().Value()
	model.Role = entity.Role()

	if err = repo.db.GormDB.Create(&model).Error; err != nil {
		return vo, err
	}

	return domain.NewAdminID(model.ID)
}

// Update Adminの更新
func (repo AdminRepository) Update(entity domain.Admin) error {
	model := Admin{}

	model.ID = entity.ID().Value()
	model.Name = entity.Name().Value()
	model.Email = entity.Email().Value()
	model.Password = entity.HashPassword().Value()
	model.Role = entity.Role()

	return repo.db.GormDB.Updates(model).Error
}

// FindByID AdminIDからAdminを取得する
func (repo AdminRepository) FindByID(id domain.AdminID) (entity domain.Admin, err error) {
	model := Admin{}

	if err = repo.db.GormDB.Where("id = ?", id.Value()).First(&model).Error; err != nil {
		return entity, err
	}

	return CreateAdminEntityFromModel(model), err
}

// FindByEmail AdminEmailからAdminを取得する
func (repo AdminRepository) FindByEmail(email domain.AdminEmail) (entity domain.Admin, err error) {
	model := Admin{}

	if err = repo.db.GormDB.Where("email = ?", email.Value()).First(&model).Error; err != nil {
		return entity, err
	}

	return CreateAdminEntityFromModel(model), err
}

// Delete AdminIDからAdminを削除する
func (repo AdminRepository) Delete(id domain.AdminID) error {
	return repo.db.GormDB.Where("id = ?", id.Value()).Delete(&Admin{}).Error
}

// ExistsByEmail AdminEmailが重複しているAdminがあるか
func (repo AdminRepository) ExistsByEmail(email domain.AdminEmail) (bool, error) {
	count := int64(0)
	err := repo.db.GormDB.Model(&Admin{}).
		Where("email = ?", email.Value()).
		Count(&count).
		Error

	return (count > 0), err
}

// ExistsByIDEmail 指定したAdminIDを除きAdminEmailが重複しているAdminがあるか
func (repo AdminRepository) ExistsByIDEmail(id domain.AdminID, email domain.AdminEmail) (bool, error) {
	count := int64(0)
	err := repo.db.GormDB.Model(&Admin{}).
		Where("id <> ? AND email = ?", id.Value(), email.Value()).
		Count(&count).
		Error

	return (count > 0), err
}

// NextIdentity 次のIDを取得する
func (repo AdminRepository) NextIdentity() (domain.AdminID, error) {
	var max uint

	sql := GetNextIdentitySelectSQL(repo.config.DB.Type)
	repo.db.GormDB.Model(&Admin{}).Select(sql).Scan(&max)

	return domain.NewAdminID(max + 1)
}

// CreateAdminEntityFromModel AdminからEntityを生成する
func CreateAdminEntityFromModel(model Admin) domain.Admin {
	return domain.NewAdmin(
		model.ID,
		model.Name,
		model.Email,
		model.Password,
		model.Role,
	)
}

// --- Admin ---

// Admin Gormモデル
type Admin struct {
	gorm.Model
	Name     string `gorm:"type:varchar(191);index;not null"`
	Email    string `gorm:"type:varchar(191);uniqueIndex;not null"`
	Password []byte
	Role     domain.AdminRole `gorm:"type:varchar(191);index;not null;default:admin"`
}
