package application

import (
	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	query "github.com/takemo101/dc-scheduler/pkg/application/query"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AppErrorType ---

const UserDuplicateError common.AppErrorType = "アカウント情報が重複しています"

// --- UserSearchInput ---

// UserSearchInput User一覧取得DTO
type UserSearchInput struct {
	Page  int
	Limit int
}

// --- UserSearchUseCase ---

// UserSearchUseCase User一覧ユースケース
type UserSearchUseCase struct {
	query query.UserQuery
}

// NewUserSearchUseCase コンストラクタ
func NewUserSearchUseCase(
	query query.UserQuery,
) UserSearchUseCase {
	return UserSearchUseCase{
		query,
	}
}

// Execute User一覧取得を実行
func (uc UserSearchUseCase) Execute(
	input UserSearchInput,
) (paginator query.UserSearchPaginatorDTO, err common.AppError) {

	parameter := query.UserSearchParameterDTO{
		Page:        input.Page,
		Limit:       input.Limit,
		OrderByKey:  "id",
		OrderByType: common.OrderByTypeDesc,
	}

	paginator, e := uc.query.Search(parameter)
	if e != nil {
		return paginator, common.NewByError(e)
	}

	return paginator, err
}

// --- UserDetailUseCase ---

// UserDetailUseCase User詳細ユースケース
type UserDetailUseCase struct {
	query query.UserQuery
}

// NewUserDetailUseCase コンストラクタ
func NewUserDetailUseCase(
	query query.UserQuery,
) UserDetailUseCase {
	return UserDetailUseCase{
		query,
	}
}

// Execute User詳細取得を実行
func (uc UserDetailUseCase) Execute(id uint) (detail query.UserDetailDTO, err common.AppError) {
	findID, e := domain.NewUserID(id)
	if e != nil {
		return detail, common.NewByError(e)
	}

	detail, e = uc.query.FindByID(findID)
	if e != nil {
		return detail, common.NewError(common.NotFoundDataError)
	}

	return detail, err
}

// --- UserStoreInput ---

// UserStoreInput User追加DTO
type UserStoreInput struct {
	Name     string
	Email    string
	Password string
	Active   bool
}

// --- UserCreateUseCase ---

// UserStoreUseCase User追加ユースケース
type UserStoreUseCase struct {
	repository domain.UserRepository
	service    domain.UserService
}

// NewUserStoreUseCase コンストラクタ
func NewUserStoreUseCase(
	repository domain.UserRepository,
	service domain.UserService,
) UserStoreUseCase {
	return UserStoreUseCase{
		repository,
		service,
	}
}

// Execute User追加を実行
func (uc UserStoreUseCase) Execute(
	input UserStoreInput,
) (id uint, err common.AppError) {

	nextID, e := uc.repository.NextIdentity()
	if e != nil {
		return id, common.NewByError(e)
	}

	entity, e := domain.CreateUser(
		nextID.Value(),
		input.Name,
		input.Email,
		input.Password,
		input.Active,
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

	return storeID.Value(), err
}

// --- UserUpdateInput ---

// UserUpdateInput User更新DTO
type UserUpdateInput struct {
	Name     string
	Email    string
	Password string
	Active   bool
}

// --- UserUpdateUseCase ---

// UserUpdateUseCase User更新ユースケース
type UserUpdateUseCase struct {
	repository domain.UserRepository
	service    domain.UserService
}

// NewUserUpdateUseCase コンストラクタ
func NewUserUpdateUseCase(
	repository domain.UserRepository,
	service domain.UserService,
) UserUpdateUseCase {
	return UserUpdateUseCase{
		repository,
		service,
	}
}

// Execute User更新を実行
func (uc UserUpdateUseCase) Execute(
	id uint,
	input UserUpdateInput,
) (err common.AppError) {
	findID, e := domain.NewUserID(id)
	if e != nil {
		return common.NewByError(e)
	}

	entity, e := uc.repository.FindByID(findID)
	if e != nil {
		return common.NewByError(e)
	}

	e = entity.Update(
		input.Name,
		input.Email,
		input.Password,
		input.Active,
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

	e = uc.repository.Update(entity)
	if e != nil {
		return common.NewByError(e)
	}

	return err
}

// --- UserDeleteUseCase ---

// UserDeleteUseCase User削除ユースケース
type UserDeleteUseCase struct {
	repository domain.UserRepository
}

// NewUserDeleteUseCase コンストラクタ
func NewUserDeleteUseCase(
	repository domain.UserRepository,
) UserDeleteUseCase {
	return UserDeleteUseCase{
		repository,
	}
}

// Execute User削除を実行
func (uc UserDeleteUseCase) Execute(id uint) (err common.AppError) {

	deleteID, e := domain.NewUserID(id)
	if e != nil {
		return common.NewByError(e)
	}

	e = uc.repository.Delete(deleteID)
	if e != nil {
		return common.NewByError(e)
	}

	return err
}
