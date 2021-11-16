package application

import (
	"time"

	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	query "github.com/takemo101/dc-scheduler/pkg/application/query"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AppErrorType ---

const RegularTimingDuplicateError common.AppErrorType = "配信タイミングが重複しています"

// --- RegularPostSearchInput ---

// RegularPostSearchInput RegularPost一覧取得DTO
type RegularPostSearchInput struct {
	Page  int
	Limit int
}

// --- RegularPostSearchUseCase ---

// RegularPostSearchUseCase RegularPost一覧ユースケース
type RegularPostSearchUseCase struct {
	query query.RegularPostQuery
}

// NewRegularPostSearchUseCase コンストラクタ
func NewRegularPostSearchUseCase(
	query query.RegularPostQuery,
) RegularPostSearchUseCase {
	return RegularPostSearchUseCase{
		query,
	}
}

// Execute RegularPost一覧取得を実行
func (uc RegularPostSearchUseCase) Execute(
	context domain.UserAuthContext,
	input RegularPostSearchInput,
) (paginator query.RegularPostSearchPaginatorDTO, err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return paginator, common.NewByError(e)
	}

	parameter := query.RegularPostSearchParameterDTO{
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

// --- RegularPostStoreInput ---

// RegularPostStoreInput RegularPost追加DTO
type RegularPostStoreInput struct {
	Message string
	Active  bool
}

// --- RegularPostStoreUseCase ---

// RegularPostStoreUseCase RegularPost追加ユースケース
type RegularPostStoreUseCase struct {
	repository    domain.RegularPostRepository
	botRepository domain.BotRepository
}

// NewRegularPostStoreUseCase コンストラクタ
func NewRegularPostStoreUseCase(
	repository domain.RegularPostRepository,
	botRepository domain.BotRepository,
) RegularPostStoreUseCase {
	return RegularPostStoreUseCase{
		repository,
		botRepository,
	}
}

// Execute RegularPost追加を実行
func (uc RegularPostStoreUseCase) Execute(
	context domain.UserAuthContext,
	botID uint,
	input RegularPostStoreInput,
) (id uint, err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return id, common.NewByError(e)
	}

	botIDVO, e := domain.NewBotID(botID)
	if e != nil {
		return id, common.NewByError(e)
	}

	bot, e := uc.botRepository.FindByID(botIDVO)
	if e != nil {
		return id, common.NewError(common.NotFoundDataError)
	}

	// ログインUserのBotかどうか
	if !bot.IsMine(auth.ID()) {
		return id, common.NewError(BotNotMineError)
	}

	nextID, e := uc.repository.NextIdentity()
	if e != nil {
		return id, common.NewByError(e)
	}

	entity, e := domain.CreateRegularPost(
		nextID.Value(),
		input.Message,
		bot,
		input.Active,
	)
	if e != nil {
		return id, common.NewByError(e)
	}

	// 情報を保存
	storeID, e := uc.repository.Store(entity)
	if e != nil {
		return id, common.NewByError(e)
	}

	return storeID.Value(), err
}

// --- RegularPostEditFormUseCase ---

// RegularPostEditFormUseCase RegularPost編集フォームユースケース
type RegularPostEditFormUseCase struct {
	repository domain.RegularPostRepository
	query      query.RegularPostQuery
}

// NewPostMessageCreateFormUseCase コンストラクタ
func NewRegularPostEditFormUseCase(
	repository domain.RegularPostRepository,
	query query.RegularPostQuery,
) RegularPostEditFormUseCase {
	return RegularPostEditFormUseCase{
		repository,
		query,
	}
}

// Execute フォーム表示のためのRegularPost取得を実行
func (uc RegularPostEditFormUseCase) Execute(
	context domain.UserAuthContext,
	id uint,
) (detail query.RegularPostDetailDTO, err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return detail, common.NewByError(e)
	}

	findID, e := domain.NewPostMessageID(id)
	if e != nil {
		return detail, common.NewByError(e)
	}

	detail, e = uc.query.FindByIDAndUserID(findID, auth.ID())
	if e != nil {
		return detail, common.NewByError(e)
	}

	return detail, err
}

// --- RegularPostUpdateInput ---

// RegularPostUpdateInput RegularPost追加DTO
type RegularPostUpdateInput struct {
	Message string
	Active  bool
}

// --- RegularPostUpdateUseCase ---

// RegularPostUpdateUseCase RegularPost追加ユースケース
type RegularPostUpdateUseCase struct {
	repository domain.RegularPostRepository
}

// NewRegularPostUpdateUseCase コンストラクタ
func NewRegularPostUpdateUseCase(
	repository domain.RegularPostRepository,
) RegularPostUpdateUseCase {
	return RegularPostUpdateUseCase{
		repository,
	}
}

// Execute RegularPost更新を実行
func (uc RegularPostUpdateUseCase) Execute(
	context domain.UserAuthContext,
	id uint,
	input RegularPostUpdateInput,
) (err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return common.NewByError(e)
	}

	idVO, e := domain.NewPostMessageID(id)
	if e != nil {
		return common.NewByError(e)
	}

	entity, e := uc.repository.FindByID(idVO)
	if e != nil {
		return common.NewError(common.NotFoundDataError)
	} else if !entity.IsMine(auth.ID()) {
		return common.NewError(BotNotMineError)
	}

	e = entity.Update(
		input.Message,
		input.Active,
	)
	if e != nil {
		return common.NewByError(e)
	}

	// 配信状態の情報を保存
	e = uc.repository.Update(entity)
	if e != nil {
		return common.NewByError(e)
	}

	return err
}

// --- RegularTimingAddInput ---

// RegularTimingInput RegularTiming追加DTO
type RegularTimingInput struct {
	DayOfWeek string
	HourTime  time.Time
}

// --- RegularTimingAddUseCase ---

// RegularTimingAddUseCase RegularTiming追加ユースケース
type RegularTimingAddUseCase struct {
	repository domain.RegularPostRepository
}

// NewRegularTimingAddUseCase コンストラクタ
func NewRegularTimingAddUseCase(
	repository domain.RegularPostRepository,
) RegularTimingAddUseCase {
	return RegularTimingAddUseCase{
		repository,
	}
}

// Execute RegularTiming追加を実行
func (uc RegularTimingAddUseCase) Execute(
	context domain.UserAuthContext,
	id uint,
	input RegularTimingInput,
) (err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return common.NewByError(e)
	}

	idVO, e := domain.NewPostMessageID(id)
	if e != nil {
		return common.NewByError(e)
	}

	entity, e := uc.repository.FindByID(idVO)
	if e != nil {
		return common.NewError(common.NotFoundDataError)
	} else if !entity.IsMine(auth.ID()) {
		return common.NewError(BotNotMineError)
	}

	e = entity.AddTiming(
		input.DayOfWeek,
		input.HourTime,
	)
	if e != nil {
		return common.NewError(RegularTimingDuplicateError)
	}

	// 更新情報を保存
	e = uc.repository.Update(entity)
	if e != nil {
		return common.NewByError(e)
	}

	return err
}

// --- RegularTimingRemoveUseCase ---

// RegularTimingRemoveUseCase RegularTiming削除ユースケース
type RegularTimingRemoveUseCase struct {
	repository domain.RegularPostRepository
}

// NewRegularTimingRemoveUseCase コンストラクタ
func NewRegularTimingRemoveUseCase(
	repository domain.RegularPostRepository,
) RegularTimingRemoveUseCase {
	return RegularTimingRemoveUseCase{
		repository,
	}
}

// Execute RegularTiming削除を実行
func (uc RegularTimingRemoveUseCase) Execute(
	context domain.UserAuthContext,
	id uint,
	input RegularTimingInput,
) (err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return common.NewByError(e)
	}

	idVO, e := domain.NewPostMessageID(id)
	if e != nil {
		return common.NewByError(e)
	}

	entity, e := uc.repository.FindByID(idVO)
	if e != nil {
		return common.NewError(common.NotFoundDataError)
	} else if !entity.IsMine(auth.ID()) {
		return common.NewError(BotNotMineError)
	}

	e = entity.RemoveTiming(
		input.DayOfWeek,
		input.HourTime,
	)
	if e != nil {
		return common.NewByError(e)
	}

	// 更新情報を保存
	e = uc.repository.Update(entity)
	if e != nil {
		return common.NewByError(e)
	}

	return err
}

// --- RegularTimingInput ---

// RegularTimingStoreInput
type RegularTimingStoreInput struct {
	DayOfWeek string
	HourTime  time.Time
}
