package application

import (
	"mime/multipart"

	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	query "github.com/takemo101/dc-scheduler/pkg/application/query"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AppErrorType ---
const (
	BotDuplicateError      common.AppErrorType = "ボット情報が重複しています"
	BotWebhookInvalidError common.AppErrorType = "ボットURLが無効です"
)

// --- BotSearchInput ---

// BotSearchInput Bot一覧取得DTO
type BotSearchInput struct {
	Page  int
	Limit int
}

// --- BotSearchUseCase ---

// BotSearchUseCase Bot一覧ユースケース
type BotSearchUseCase struct {
	query query.BotQuery
}

// NewBotSearchUseCase コンストラクタ
func NewBotSearchUseCase(
	query query.BotQuery,
) BotSearchUseCase {
	return BotSearchUseCase{
		query,
	}
}

// Execute Bot一覧取得を実行
func (uc BotSearchUseCase) Execute(
	context domain.UserAuthContext,
	input BotSearchInput,
) (paginator query.BotSearchPaginatorDTO, err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return paginator, common.NewByError(e)
	}

	parameter := query.BotSearchParameterDTO{
		Page:        input.Page,
		Limit:       input.Limit,
		OrderByKey:  "id",
		OrderByType: common.OrderByTypeDesc,
	}

	paginator, e = uc.query.SearchByUserID(parameter, auth.ID())
	if e != nil {
		return paginator, common.NewByError(e)
	}

	return paginator, err
}

// --- BotDetailUseCase ---

// BotDetailUseCase Bot詳細ユースケース
type BotDetailUseCase struct {
	repository domain.BotRepository
	query      query.BotQuery
}

// NewBotDetailUseCase コンストラクタ
func NewBotDetailUseCase(
	repository domain.BotRepository,
	query query.BotQuery,
) BotDetailUseCase {
	return BotDetailUseCase{
		repository,
		query,
	}
}

// Execute Bot詳細取得を実行
func (uc BotDetailUseCase) Execute(
	context domain.UserAuthContext,
	id uint,
) (detail query.BotDetailDTO, err common.AppError) {
	findID, e := domain.NewBotID(id)
	if e != nil {
		return detail, common.NewByError(e)
	}

	entity, e := uc.repository.FindByID(findID)
	if e != nil {
		return detail, common.NewError(common.NotFoundDataError)
	}

	// ポリシーチェック
	policy := domain.NewUserBotPolicy(context)
	ok, e := policy.Detail(entity)
	if e != nil {
		return detail, common.NewByError(e)
	} else if !ok {
		return detail, common.NewError(common.NotTargetOwnerError)
	}

	detail, e = uc.query.FindByID(findID)
	if e != nil {
		return detail, common.NewError(common.NotFoundDataError)
	}

	return detail, err
}

// --- BotStoreInput ---

// BotStoreInput Bot追加DTO
type BotStoreInput struct {
	Name            string
	AtatarFile      *multipart.FileHeader
	AtatarDirectory string
	Webhook         string
	Active          bool
}

// --- BotCreateUseCase ---

// BotStoreUseCase Bot追加ユースケース
type BotStoreUseCase struct {
	repository     domain.BotRepository
	adapter        domain.DiscordWebhookCheckAdapter
	fileRepository domain.BotAtatarImageRepository
	service        domain.BotService
}

// NewBotStoreUseCase コンストラクタ
func NewBotStoreUseCase(
	repository domain.BotRepository,
	adapter domain.DiscordWebhookCheckAdapter,
	fileRepository domain.BotAtatarImageRepository,
	service domain.BotService,
) BotStoreUseCase {
	return BotStoreUseCase{
		repository,
		adapter,
		fileRepository,
		service,
	}
}

// Execute Bot追加を実行
func (uc BotStoreUseCase) Execute(
	context domain.UserAuthContext,
	input BotStoreInput,
) (id uint, err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return id, common.NewByError(e)
	}

	var avatar string

	// アバターがアップロードされている場合
	if input.AtatarFile != nil {
		avatarEntity, e := domain.UploadBotAtatarImage(
			input.AtatarFile,
			input.AtatarDirectory,
		)
		if e != nil {
			return id, common.NewByError(e)
		}

		avatarVO, e := uc.fileRepository.Store(avatarEntity)
		avatar = avatarVO.Value()
	}

	nextID, e := uc.repository.NextIdentity()
	if e != nil {
		return id, common.NewByError(e)
	}

	entity, e := domain.CreateBot(
		nextID.Value(),
		input.Name,
		avatar,
		input.Webhook,
		input.Active,
		auth.ID(),
	)
	if e != nil {
		return id, common.NewByError(e)
	}

	// ウェブフックの重複チェック
	duplicate, e := uc.service.IsDuplicate(entity)
	if e != nil {
		return id, common.NewByError(e)
	}
	if duplicate {
		return id, common.NewError(BotDuplicateError)
	}

	// ウェブフックの有効性チェック
	ok, _ := uc.adapter.Check(entity.Webhook())
	if !ok {
		return id, common.NewError(BotWebhookInvalidError)
	}

	storeID, e := uc.repository.Store(entity)

	if e != nil {
		return id, common.NewByError(e)
	}

	return storeID.Value(), err
}

// --- BotUpdateInput ---

// BotUpdateInput Bot更新DTO
type BotUpdateInput struct {
	Name            string
	AtatarFile      *multipart.FileHeader
	AtatarDirectory string
	Webhook         string
	Active          bool
}

// --- BotUpdateUseCase ---

// BotUpdateUseCase Bot更新ユースケース
type BotUpdateUseCase struct {
	repository     domain.BotRepository
	adapter        domain.DiscordWebhookCheckAdapter
	fileRepository domain.BotAtatarImageRepository
	service        domain.BotService
}

// NewBotUpdateUseCase コンストラクタ
func NewBotUpdateUseCase(
	repository domain.BotRepository,
	adapter domain.DiscordWebhookCheckAdapter,
	fileRepository domain.BotAtatarImageRepository,
	service domain.BotService,
) BotUpdateUseCase {
	return BotUpdateUseCase{
		repository,
		adapter,
		fileRepository,
		service,
	}
}

// Execute Bot更新を実行
func (uc BotUpdateUseCase) Execute(
	context domain.UserAuthContext,
	id uint,
	input BotUpdateInput,
) (err common.AppError) {
	findID, e := domain.NewBotID(id)
	if e != nil {
		return common.NewByError(e)
	}

	entity, e := uc.repository.FindByID(findID)
	if e != nil {
		return common.NewByError(e)
	}

	// ポリシーチェック
	policy := domain.NewUserBotPolicy(context)
	ok, e := policy.Update(entity)
	if e != nil {
		return common.NewByError(e)
	} else if !ok {
		return common.NewError(common.NotTargetOwnerError)
	}

	var avatar string

	// アバターがアップロードされている場合
	if input.AtatarFile != nil {
		avatarEntity, e := domain.UploadBotAtatarImage(
			input.AtatarFile,
			input.AtatarDirectory,
		)
		if e != nil {
			return common.NewByError(e)
		}

		avatarVO, e := uc.fileRepository.Update(
			avatarEntity,
			entity.Atatar(),
		)
		avatar = avatarVO.Value()
	} else {
		avatar = entity.Atatar().Value()
	}

	e = entity.Update(
		input.Name,
		avatar,
		input.Webhook,
		input.Active,
	)
	if e != nil {
		return common.NewByError(e)
	}

	// ウェブフックの重複チェック
	duplicate, e := uc.service.IsDuplicateWithoutSelf(entity)
	if e != nil {
		return common.NewByError(e)
	}
	if duplicate {
		return common.NewError(BotDuplicateError)
	}

	// ウェブフックの有効性チェック
	ok, _ = uc.adapter.Check(entity.Webhook())
	if !ok {
		return common.NewError(BotWebhookInvalidError)
	}

	e = uc.repository.Update(entity)
	if e != nil {
		return common.NewByError(e)
	}

	return err
}

// --- BotDeleteUseCase ---

// BotDeleteUseCase Bot削除ユースケース
type BotDeleteUseCase struct {
	repository     domain.BotRepository
	fileRepository domain.BotAtatarImageRepository
}

// NewBotDeleteUseCase コンストラクタ
func NewBotDeleteUseCase(
	repository domain.BotRepository,
	fileRepository domain.BotAtatarImageRepository,
) BotDeleteUseCase {
	return BotDeleteUseCase{
		repository,
		fileRepository,
	}
}

// Execute Bot削除を実行
func (uc BotDeleteUseCase) Execute(
	context domain.UserAuthContext,
	id uint,
) (err common.AppError) {
	deleteID, e := domain.NewBotID(id)
	if e != nil {
		return common.NewByError(e)
	}

	entity, e := uc.repository.FindByID(deleteID)
	if e != nil {
		return common.NewByError(e)
	}

	// ポリシーチェック
	policy := domain.NewUserBotPolicy(context)
	ok, e := policy.Delete(entity)
	if e != nil {
		return common.NewByError(e)
	} else if !ok {
		return common.NewError(common.NotTargetOwnerError)
	}

	e = uc.repository.Delete(deleteID)
	if e != nil {
		return common.NewByError(e)
	}

	e = uc.fileRepository.Delete(entity.Atatar())
	if e != nil {
		return common.NewByError(e)
	}

	return err
}
