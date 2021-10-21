package application

import (
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AppErrorType ---

const AdminDuplicateError AppErrorType = "アカウント情報が重複しています"

// --- AdminSearchInput ---

// AdminSearchInput Admin一覧取得DTO
type AdminSearchInput struct {
	Page  int
	Limit int
}

// --- AdminSearchUseCase ---

// AdminSearchUseCase Admin一覧ユースケース
type AdminSearchUseCase struct {
	query AdminQuery
}

// NewAdminSearchUseCase コンストラクタ
func NewAdminSearchUseCase(
	query AdminQuery,
) AdminSearchUseCase {
	return AdminSearchUseCase{
		query,
	}
}

// Execute Admin一覧取得を実行
func (uc AdminSearchUseCase) Execute(
	input AdminSearchInput,
) (paginator AdminSearchPaginatorDTO, err AppError) {

	parameter := AdminSearchParameterDTO{
		Page:        input.Page,
		Limit:       input.Limit,
		OrderByKey:  "id",
		OrderByType: OrderByTypeDesc,
	}

	paginator, e := uc.query.Search(parameter)
	if e != nil {
		return paginator, NewByError(e)
	}

	return paginator, err
}

// --- AdminDetailUseCase ---

// AdminDetailUseCase Admin詳細ユースケース
type AdminDetailUseCase struct {
	query AdminQuery
}

// NewAdminDetailUseCase コンストラクタ
func NewAdminDetailUseCase(
	query AdminQuery,
) AdminDetailUseCase {
	return AdminDetailUseCase{
		query,
	}
}

// Execute Admin詳細取得を実行
func (uc AdminDetailUseCase) Execute(id uint) (detail AdminDetailDTO, err AppError) {
	findID, e := domain.NewAdminID(id)
	if e != nil {
		return detail, NewByError(e)
	}

	detail, e = uc.query.FindByID(findID)
	if e != nil {
		return detail, NewError(NotFoundDataError)
	}

	return detail, err
}

// --- AdminStoreInput ---

// AdminStoreInput Admin追加DTO
type AdminStoreInput struct {
	Name     string
	Email    string
	Role     string
	Password string
}

// --- AdminCreateUseCase ---

// AdminStoreUseCase Admin追加ユースケース
type AdminStoreUseCase struct {
	repository domain.AdminRepository
	service    domain.AdminService
}

// NewAdminStoreUseCase コンストラクタ
func NewAdminStoreUseCase(
	repository domain.AdminRepository,
	service domain.AdminService,
) AdminStoreUseCase {
	return AdminStoreUseCase{
		repository,
		service,
	}
}

// Execute Admin追加を実行
func (uc AdminStoreUseCase) Execute(
	input AdminStoreInput,
) (id uint, err AppError) {

	nextID, e := uc.repository.NextIdentity()
	if e != nil {
		return id, NewByError(e)
	}

	entity, e := domain.CreateAdmin(
		nextID.Value(),
		input.Name,
		input.Email,
		input.Password,
		input.Role,
	)
	if e != nil {
		return id, NewByError(e)
	}

	// メールアドレスの重複チェック
	duplicate, e := uc.service.IsDuplicate(entity)
	if e != nil {
		return id, NewByError(e)
	}
	if duplicate {
		return id, NewError(AdminDuplicateError)
	}

	storeID, e := uc.repository.Store(entity)
	if e != nil {
		return id, NewByError(e)
	}

	return storeID.Value(), err
}

// --- AdminUpdateInput ---

// AdminUpdateInput Admin更新DTO
type AdminUpdateInput struct {
	Name     string
	Email    string
	Role     string
	Password string
}

// --- AdminUpdateUseCase ---

// AdminUpdateUseCase Admin更新ユースケース
type AdminUpdateUseCase struct {
	repository domain.AdminRepository
	service    domain.AdminService
}

// NewAdminUpdateUseCase コンストラクタ
func NewAdminUpdateUseCase(
	repository domain.AdminRepository,
	service domain.AdminService,
) AdminUpdateUseCase {
	return AdminUpdateUseCase{
		repository,
		service,
	}
}

// Execute Admin更新を実行
func (uc AdminUpdateUseCase) Execute(
	id uint,
	input AdminUpdateInput,
) (err AppError) {
	findID, e := domain.NewAdminID(id)
	if e != nil {
		return NewByError(e)
	}

	entity, e := uc.repository.FindByID(findID)
	if e != nil {
		return NewByError(e)
	}

	e = entity.Update(
		input.Name,
		input.Email,
		input.Password,
		input.Role,
	)
	if e != nil {
		return NewByError(e)
	}

	duplicate, e := uc.service.IsDuplicateWithoutSelf(entity)
	if e != nil {
		return NewByError(e)
	}
	if duplicate {
		return NewError(AdminDuplicateError)
	}

	e = uc.repository.Update(entity)
	if e != nil {
		return NewByError(e)
	}

	return err
}

// --- AdminDeleteUseCase ---

// AdminDeleteUseCase Admin削除ユースケース
type AdminDeleteUseCase struct {
	repository domain.AdminRepository
}

// NewAdminDeleteUseCase コンストラクタ
func NewAdminDeleteUseCase(
	repository domain.AdminRepository,
) AdminDeleteUseCase {
	return AdminDeleteUseCase{
		repository,
	}
}

// Execute Admin削除を実行
func (uc AdminDeleteUseCase) Execute(id uint) (err AppError) {

	deleteID, e := domain.NewAdminID(id)
	if e != nil {
		return NewByError(e)
	}

	e = uc.repository.Delete(deleteID)
	if e != nil {
		return NewByError(e)
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
	query AdminQuery
}

// NewMyAccountDetailUseCase コンストラクタ
func NewMyAccountDetailUseCase(
	query AdminQuery,
) MyAccountDetailUseCase {
	return MyAccountDetailUseCase{
		query,
	}
}

// Execute アカウント詳細取得を実行
func (uc MyAccountDetailUseCase) Execute(
	context domain.AdminAuthContext,
) (detail AdminDetailDTO, err AppError) {

	auth, e := context.AdminAuth()
	if e != nil {
		return detail, NewByError(e)
	}

	detail, e = uc.query.FindByID(auth.ID())
	if e != nil {
		return detail, NewByError(e)
	}

	return detail, err
}

// --- MyAccountUpdateUseCase ---

// MyAccountUpdateUseCase アカウント更新ユースケース
type MyAccountUpdateUseCase struct {
	repository domain.AdminRepository
	service    domain.AdminService
}

// NewMyAccountUpdateUseCase コンストラクタ
func NewMyAccountUpdateUseCase(
	repository domain.AdminRepository,
	service domain.AdminService,
) MyAccountUpdateUseCase {
	return MyAccountUpdateUseCase{
		repository,
		service,
	}
}

// Execute アカウント更新を実行
func (uc MyAccountUpdateUseCase) Execute(
	context domain.AdminAuthContext,
	input MyAccountUpdateInput,
) (err AppError) {

	auth, e := context.AdminAuth()
	if e != nil {
		return NewByError(e)
	}

	entity, e := uc.repository.FindByID(auth.ID())
	if e != nil {
		return NewByError(e)
	}

	e = entity.Update(
		input.Name,
		input.Email,
		input.Password,
		auth.Role().Value(),
	)
	if e != nil {
		return NewByError(e)
	}

	duplicate, e := uc.service.IsDuplicateWithoutSelf(entity)
	if e != nil {
		return NewByError(e)
	}
	if duplicate {
		return NewError(AdminDuplicateError)
	}

	admin := entity.CreateLoginAuth()

	// DBへ永続化
	e = uc.repository.Update(entity)
	if e != nil {
		return NewByError(e)
	}

	// セッションへ永続化
	e = context.Login(admin)
	if e != nil {
		return NewByError(e)
	}

	return err
}
