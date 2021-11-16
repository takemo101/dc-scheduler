package application

import (
	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AppErrorType ---

const UserNotFoundAccountError common.AppErrorType = "アカウントが見つかりません"

// --- UserLoginInput ---

// UserLoginInput ログインリクエストDTO
type UserLoginInput struct {
	Email    string
	Password string
}

// --- UserLoginUseCase ---

// UserLoginUseCase ログインユースケース
type UserLoginUseCase struct {
	repository domain.UserRepository
}

// NewUserLoginUseCase コンストラクタ
func NewUserLoginUseCase(
	repository domain.UserRepository,
) UserLoginUseCase {
	return UserLoginUseCase{
		repository,
	}
}

// Execute ログインを実行
func (uc UserLoginUseCase) Execute(
	context domain.UserAuthContext,
	input UserLoginInput,
) (err common.AppError) {

	loginError := common.NewError(UserNotFoundAccountError)

	// メールアドレスからUserを取得
	entity, e := uc.repository.FindByEmail(domain.UserEmail(input.Email))
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

	auth, e := entity.CreateLoginAuth()
	if e != nil {
		return loginError
	}

	// ログイン
	e = context.Login(auth)
	if e != nil {
		return common.NewByError(e)
	}

	return err
}

// --- UserLogoutUseCase ---

// UserLooutUseCase ログアウトユースケース
type UserLogoutUseCase struct{}

// NewUserLogoutUseCase コンストラクタ
func NewUserLogoutUseCase() UserLogoutUseCase {
	return UserLogoutUseCase{}
}

// Execute ログアウトを実行
func (uc UserLogoutUseCase) Execute(context domain.UserAuthContext) (err common.AppError) {
	e := context.Logout()
	if e != nil {
		return common.NewByError(e)
	}

	return err
}
