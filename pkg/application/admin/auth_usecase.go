package application

import (
	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AppErrorType ---

const AdminNotFoundAccountError common.AppErrorType = "アカウントが見つかりません"

// --- AdminLoginInput ---

// AdminLoginInput ログインリクエストDTO
type AdminLoginInput struct {
	Email    string
	Password string
}

// --- AdminLoginUseCase ---

// AdminLoginUseCase ログインユースケース
type AdminLoginUseCase struct {
	repository domain.AdminRepository
}

// NewAdminLoginUseCase コンストラクタ
func NewAdminLoginUseCase(
	repository domain.AdminRepository,
) AdminLoginUseCase {
	return AdminLoginUseCase{
		repository,
	}
}

// Execute ログインを実行
func (uc AdminLoginUseCase) Execute(
	context domain.AdminAuthContext,
	input AdminLoginInput,
) (err common.AppError) {

	loginError := common.NewError(AdminNotFoundAccountError)

	// メールアドレスからAdminを取得
	entity, e := uc.repository.FindByEmail(domain.AdminEmail(input.Email))
	if e != nil {
		return loginError
	}

	// パスワードチェック
	ok, e := entity.ComparePassword(input.Password)
	if e != nil {
		return common.NewByError(e)
	}

	if !ok {
		return loginError
	}

	// ログイン
	e = context.Login(entity.CreateLoginAuth())
	if e != nil {
		return common.NewByError(e)
	}

	return err
}

// --- AdminLogoutUseCase ---

// AdminLooutUseCase ログアウトユースケース
type AdminLogoutUseCase struct{}

// NewAdminLogoutUseCase コンストラクタ
func NewAdminLogoutUseCase() AdminLogoutUseCase {
	return AdminLogoutUseCase{}
}

// Execute ログアウトを実行
func (uc AdminLogoutUseCase) Execute(context domain.AdminAuthContext) (err common.AppError) {
	e := context.Logout()
	if e != nil {
		return common.NewByError(e)
	}

	return err
}
