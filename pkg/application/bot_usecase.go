package application

import (
	"mime/multipart"

	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AppErrorType ---

const BotDuplicateError AppErrorType = "ボット情報が重複しています"

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
		Page:    input.Page,
		Limit:   input.Limit,
		OrderBy: "id DESC",
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
		return detail, NewError(NotFoundDataError)
	}

	detail, e = uc.query.FindByID(findID)
	if e != nil {
		return detail, NewByError(e)
	}

	return detail, err
}

// --- BotStoreInput ---

// BotStoreInput Bot追加DTO
type BotStoreInput struct {
	Name            string
	AvatorFile      *multipart.FileHeader
	AvatorDirectory string
	Webhook         string
	Active          bool
}

// --- BotCreateUseCase ---

// BotStoreUseCase Bot追加ユースケース
type BotStoreUseCase struct {
	repository     domain.BotRepository
	fileRepository domain.BotAvatorImageRepository
	service        domain.BotService
}

// NewBotStoreUseCase コンストラクタ
func NewBotStoreUseCase(
	repository domain.BotRepository,
	fileRepository domain.BotAvatorImageRepository,
	service domain.BotService,
) BotStoreUseCase {
	return BotStoreUseCase{
		repository,
		fileRepository,
		service,
	}
}

// Execute Bot追加を実行
func (uc BotStoreUseCase) Execute(
	input BotStoreInput,
) (id uint, err AppError) {
	var avator string

	// アバターがアップロードされている場合
	if input.AvatorFile != nil {
		avatorEntity, e := domain.UploadBotAvatorImage(
			input.AvatorFile,
			input.AvatorDirectory,
		)
		if e != nil {
			return id, NewByError(e)
		}

		avatorVO, e := uc.fileRepository.Store(avatorEntity)
		avator = avatorVO.Value()
	}

	nextID, e := uc.repository.NextIdentity()
	if e != nil {
		return id, NewByError(e)
	}

	entity, e := domain.CreateBot(
		nextID.Value(),
		input.Name,
		avator,
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
	AvatorFile      *multipart.FileHeader
	AvatorDirectory string
	Webhook         string
	Active          bool
}

// --- BotUpdateUseCase ---

// BotUpdateUseCase Bot更新ユースケース
type BotUpdateUseCase struct {
	repository     domain.BotRepository
	fileRepository domain.BotAvatorImageRepository
	service        domain.BotService
}

// NewBotUpdateUseCase コンストラクタ
func NewBotUpdateUseCase(
	repository domain.BotRepository,
	fileRepository domain.BotAvatorImageRepository,
	service domain.BotService,
) BotUpdateUseCase {
	return BotUpdateUseCase{
		repository,
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

	var avator string

	// アバターがアップロードされている場合
	if input.AvatorFile != nil {
		avatorEntity, e := domain.UploadBotAvatorImage(
			input.AvatorFile,
			input.AvatorDirectory,
		)
		if e != nil {
			return NewByError(e)
		}

		avatorVO, e := uc.fileRepository.Update(
			avatorEntity,
			entity.Avator(),
		)
		avator = avatorVO.Value()
	}

	e = entity.Update(
		input.Name,
		avator,
		input.Webhook,
		input.Active,
	)
	if e != nil {
		return NewByError(e)
	}

	duplicate, e := uc.service.IsDuplicateWithoutSelf(entity)
	if e != nil {
		return NewByError(e)
	}
	if duplicate {
		return NewError(BotDuplicateError)
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
	fileRepository domain.BotAvatorImageRepository
}

// NewBotDeleteUseCase コンストラクタ
func NewBotDeleteUseCase(
	repository domain.BotRepository,
	fileRepository domain.BotAvatorImageRepository,
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

	e = uc.fileRepository.Delete(entity.Avator())
	if e != nil {
		return NewByError(e)
	}

	return err
}
