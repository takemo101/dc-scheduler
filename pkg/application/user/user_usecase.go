package application

import (
	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	query "github.com/takemo101/dc-scheduler/pkg/application/query"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AppErrorType ---

const (
	UserDuplicateError             common.AppErrorType = "アカウント情報が重複しています"
	UserNotMatchActivationKeyError common.AppErrorType = "アクティベーションキーが一致しません"
)

// --- UserRegistInput ---

// UserRegistInput User仮登録DTO
type UserRegistInput struct {
	Name     string
	Email    string
	Password string
}

// --- UserRegistUseCase ---

// UserRegistUseCase User追加ユースケース
type UserRegistUseCase struct {
	repository domain.UserRepository
	service    domain.UserService
	crypter    domain.UserActivationSignatureCrypter
	sender     domain.TemplateMailSender
}

// NewUserRegistUseCase コンストラクタ
func NewUserRegistUseCase(
	repository domain.UserRepository,
	service domain.UserService,
	crypter domain.UserActivationSignatureCrypter,
	sender domain.TemplateMailSender,
) UserRegistUseCase {
	return UserRegistUseCase{
		repository,
		service,
		crypter,
		sender,
	}
}

// Execute User仮登録を実行
func (uc UserRegistUseCase) Execute(
	input UserRegistInput,
) (id uint, err common.AppError) {

	nextID, e := uc.repository.NextIdentity()
	if e != nil {
		return id, common.NewByError(e)
	}

	entity, e := domain.CreateTemporaryUser(
		nextID.Value(),
		input.Name,
		input.Email,
		input.Password,
	)
	if e != nil {
		return id, common.NewByError(e)
	}

	// メールアドレスの重複チェック
	duplicate, e := uc.service.IsDuplicate(entity)
	if e != nil {
		return id, common.NewByError(e)
	}
	if duplicate {
		return id, common.NewError(UserDuplicateError)
	}

	storeID, e := uc.repository.Store(entity)
	if e != nil {
		return id, common.NewByError(e)
	}

	// 暗号化してシグネチャを生成する
	signature, e := uc.crypter.Encrypt(entity.CreateSignature())
	if e != nil {
		return storeID.Value(), common.NewByError(e)
	}

	// 通知情報を生成
	notify, e := domain.CreateUserActivationNotify(entity, signature)
	if e != nil {
		return storeID.Value(), common.NewByError(e)
	}

	// 通知する
	e = uc.sender.Send(notify)
	if e != nil {
		return storeID.Value(), common.NewByError(e)
	}

	return storeID.Value(), err
}

// --- UserActivationUseCase ---

// UserActivationUseCase Userアクティベーションユースケース
type UserActivationUseCase struct {
	repository domain.UserRepository
	crypter    domain.UserActivationSignatureCrypter
}

// NewUserActivationUseCase コンストラクタ
func NewUserActivationUseCase(
	repository domain.UserRepository,
	crypter domain.UserActivationSignatureCrypter,
) UserActivationUseCase {
	return UserActivationUseCase{
		repository,
		crypter,
	}
}

// Execute Userアクティベーションの実行
func (uc UserActivationUseCase) Execute(
	signature string,
) (err common.AppError) {

	// シグネチャを複合
	userSignature, e := uc.crypter.Decrypt(signature)
	if e != nil {
		return common.NewByError(e)
	}

	entity, e := uc.repository.FindByEmail(userSignature.Email)
	if e != nil {
		return common.NewByError(e)
	}

	// アクティベーション
	ok, e := entity.Activation(userSignature.ActivationKey)
	if e != nil {
		return common.NewByError(e)
	}
	if !ok {
		return common.NewError(UserNotMatchActivationKeyError)
	}

	if e = uc.repository.Update(entity); e != nil {
		return common.NewByError(e)
	}

	return err
}

// --- MyAccountUpdateInput ---

// MyAccountUpdateInput アカウント更新DTO
type MyAccountUpdateInput struct {
	Name     string
	Email    string
	Password string
}

// --- MyAccountDetailUseCase ---

// MyAccountDetailUseCase アカウント取得ユースケース
type MyAccountDetailUseCase struct {
	query query.UserQuery
}

// NewMyAccountDetailUseCase コンストラクタ
func NewMyAccountDetailUseCase(
	query query.UserQuery,
) MyAccountDetailUseCase {
	return MyAccountDetailUseCase{
		query,
	}
}

// Execute アカウント詳細取得を実行
func (uc MyAccountDetailUseCase) Execute(
	context domain.UserAuthContext,
) (detail query.UserDetailDTO, err common.AppError) {

	auth, e := context.UserAuth()
	if e != nil {
		return detail, common.NewByError(e)
	}

	detail, e = uc.query.FindByID(auth.ID())
	if e != nil {
		return detail, common.NewByError(e)
	}

	return detail, err
}

// --- MyAccountUpdateUseCase ---

// MyAccountUpdateUseCase アカウント更新ユースケース
type MyAccountUpdateUseCase struct {
	repository domain.UserRepository
	service    domain.UserService
}

// NewMyAccountUpdateUseCase コンストラクタ
func NewMyAccountUpdateUseCase(
	repository domain.UserRepository,
	service domain.UserService,
) MyAccountUpdateUseCase {
	return MyAccountUpdateUseCase{
		repository,
		service,
	}
}

// Execute アカウント更新を実行
func (uc MyAccountUpdateUseCase) Execute(
	context domain.UserAuthContext,
	input MyAccountUpdateInput,
) (err common.AppError) {

	auth, e := context.UserAuth()
	if e != nil {
		return common.NewByError(e)
	}

	entity, e := uc.repository.FindByID(auth.ID())
	if e != nil {
		return common.NewByError(e)
	}

	e = entity.Update(
		input.Name,
		input.Email,
		input.Password,
		true,
	)
	if e != nil {
		return common.NewByError(e)
	}

	duplicate, e := uc.service.IsDuplicateWithoutSelf(entity)
	if e != nil {
		return common.NewByError(e)
	}
	if duplicate {
		return common.NewError(UserDuplicateError)
	}

	user, e := entity.CreateLoginAuth()
	if e != nil {
		return common.NewByError(e)
	}

	// DBへ永続化
	e = uc.repository.Update(entity)
	if e != nil {
		return common.NewByError(e)
	}

	// セッションへ永続化
	e = context.Login(user)
	if e != nil {
		return common.NewByError(e)
	}

	return err
}
