package application

import (
	"mime/multipart"

	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AppErrorType ---

const BotDuplicateError AppErrorType = "ボット情報が重複しています"
const BotWebhookInvalidError AppErrorType = "ボットURLが無効です"

// --- BotSearchInput ---

// BotSearchInput Bot一覧取得DTO
type BotSearchInput struct {
	Page  int
	Limit int
}

// --- BotSearchUseCase ---

// BotSearchUseCase Bot一覧ユースケース
type BotSearchUseCase struct {
	query BotQuery
}

// NewBotSearchUseCase コンストラクタ
func NewBotSearchUseCase(
	query BotQuery,
) BotSearchUseCase {
	return BotSearchUseCase{
		query,
	}
}

// Execute Bot一覧取得を実行
func (uc BotSearchUseCase) Execute(
	input BotSearchInput,
) (paginator BotSearchPaginatorDTO, err AppError) {

	parameter := BotSearchParameterDTO{
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

// --- BotDetailUseCase ---

// BotDetailUseCase Bot詳細ユースケース
type BotDetailUseCase struct {
	query BotQuery
}

// NewBotDetailUseCase コンストラクタ
func NewBotDetailUseCase(
	query BotQuery,
) BotDetailUseCase {
	return BotDetailUseCase{
		query,
	}
}

// Execute Bot詳細取得を実行
func (uc BotDetailUseCase) Execute(id uint) (detail BotDetailDTO, err AppError) {
	findID, e := domain.NewBotID(id)
	if e != nil {
		return detail, NewByError(e)
	}

	detail, e = uc.query.FindByID(findID)
	if e != nil {
		return detail, NewError(NotFoundDataError)
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
	input BotStoreInput,
) (id uint, err AppError) {
	var avatar string

	// アバターがアップロードされている場合
	if input.AtatarFile != nil {
		avatarEntity, e := domain.UploadBotAtatarImage(
			input.AtatarFile,
			input.AtatarDirectory,
		)
		if e != nil {
			return id, NewByError(e)
		}

		avatarVO, e := uc.fileRepository.Store(avatarEntity)
		avatar = avatarVO.Value()
	}

	nextID, e := uc.repository.NextIdentity()
	if e != nil {
		return id, NewByError(e)
	}

	entity, e := domain.CreateBot(
		nextID.Value(),
		input.Name,
		avatar,
		input.Webhook,
		input.Active,
	)
	if e != nil {
		return id, NewByError(e)
	}

	// ウェブフックの重複チェック
	duplicate, e := uc.service.IsDuplicate(entity)
	if e != nil {
		return id, NewByError(e)
	}
	if duplicate {
		return id, NewError(BotDuplicateError)
	}

	// ウェブフックの有効性チェック
	ok, _ := uc.adapter.Check(entity.Webhook())
	if !ok {
		return id, NewError(BotWebhookInvalidError)
	}

	storeID, e := uc.repository.Store(entity)

	if e != nil {
		return id, NewByError(e)
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
	id uint,
	input BotUpdateInput,
) (err AppError) {
	findID, e := domain.NewBotID(id)
	if e != nil {
		return NewByError(e)
	}

	entity, e := uc.repository.FindByID(findID)
	if e != nil {
		return NewByError(e)
	}

	var avatar string

	// アバターがアップロードされている場合
	if input.AtatarFile != nil {
		avatarEntity, e := domain.UploadBotAtatarImage(
			input.AtatarFile,
			input.AtatarDirectory,
		)
		if e != nil {
			return NewByError(e)
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
		return NewByError(e)
	}

	// ウェブフックの重複チェック
	duplicate, e := uc.service.IsDuplicateWithoutSelf(entity)
	if e != nil {
		return NewByError(e)
	}
	if duplicate {
		return NewError(BotDuplicateError)
	}

	// ウェブフックの有効性チェック
	ok, _ := uc.adapter.Check(entity.Webhook())
	if !ok {
		return NewError(BotWebhookInvalidError)
	}

	e = uc.repository.Update(entity)
	if e != nil {
		return NewByError(e)
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
func (uc BotDeleteUseCase) Execute(id uint) (err AppError) {

	deleteID, e := domain.NewBotID(id)
	if e != nil {
		return NewByError(e)
	}

	entity, e := uc.repository.FindByID(deleteID)
	if e != nil {
		return NewByError(e)
	}

	e = uc.repository.Delete(deleteID)
	if e != nil {
		return NewByError(e)
	}

	e = uc.fileRepository.Delete(entity.Atatar())
	if e != nil {
		return NewByError(e)
	}

	return err
}
