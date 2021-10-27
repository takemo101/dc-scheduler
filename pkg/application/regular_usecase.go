package application

import (
	"time"

	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AppErrorType ---

const RegularTimingDuplicateError AppErrorType = "配信タイミングが重複しています"

// --- RegularPostSearchInput ---

// RegularPostSearchInput RegularPost一覧取得DTO
type RegularPostSearchInput struct {
	Page  int
	Limit int
}

// --- RegularPostSearchUseCase ---

// RegularPostSearchUseCase RegularPost一覧ユースケース
type RegularPostSearchUseCase struct {
	query RegularPostQuery
}

// NewRegularPostSearchUseCase コンストラクタ
func NewRegularPostSearchUseCase(
	query RegularPostQuery,
) RegularPostSearchUseCase {
	return RegularPostSearchUseCase{
		query,
	}
}

// Execute RegularPost一覧取得を実行
func (uc RegularPostSearchUseCase) Execute(
	input RegularPostSearchInput,
) (paginator RegularPostSearchPaginatorDTO, err AppError) {

	parameter := RegularPostSearchParameterDTO{
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
	botID uint,
	input RegularPostStoreInput,
) (id uint, err AppError) {

	botIDVO, e := domain.NewBotID(botID)
	if e != nil {
		return id, NewByError(e)
	}

	bot, e := uc.botRepository.FindByID(botIDVO)
	if e != nil {
		return id, NewError(NotFoundDataError)
	}

	nextID, e := uc.repository.NextIdentity()
	if e != nil {
		return id, NewByError(e)
	}

	entity, e := domain.CreateRegularPost(
		nextID.Value(),
		input.Message,
		bot,
		input.Active,
	)
	if e != nil {
		return id, NewByError(e)
	}

	// 情報を保存
	storeID, e := uc.repository.Store(entity)
	if e != nil {
		return id, NewByError(e)
	}

	return storeID.Value(), err
}

// --- RegularPostEditFormUseCase ---

// RegularPostEditFormUseCase RegularPost編集フォームユースケース
type RegularPostEditFormUseCase struct {
	repository domain.RegularPostRepository
	query      RegularPostQuery
}

// NewPostMessageCreateFormUseCase コンストラクタ
func NewRegularPostEditFormUseCase(
	repository domain.RegularPostRepository,
	query RegularPostQuery,
) RegularPostEditFormUseCase {
	return RegularPostEditFormUseCase{
		repository,
		query,
	}
}

// Execute フォーム表示のためのRegularPost取得を実行
func (uc RegularPostEditFormUseCase) Execute(id uint) (detail RegularPostDetailDTO, err AppError) {
	findID, e := domain.NewPostMessageID(id)
	if e != nil {
		return detail, NewByError(e)
	}

	detail, e = uc.query.FindByID(findID)
	if e != nil {
		return detail, NewByError(e)
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
	id uint,
	input RegularPostUpdateInput,
) (err AppError) {

	idVO, e := domain.NewPostMessageID(id)
	if e != nil {
		return NewByError(e)
	}

	entity, e := uc.repository.FindByID(idVO)
	if e != nil {
		return NewError(NotFoundDataError)
	}

	e = entity.Update(
		input.Message,
		input.Active,
	)
	if e != nil {
		return NewByError(e)
	}

	// 配信状態の情報を保存
	e = uc.repository.Update(entity)
	if e != nil {
		return NewByError(e)
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
	id uint,
	input RegularTimingInput,
) (err AppError) {

	idVO, e := domain.NewPostMessageID(id)
	if e != nil {
		return NewByError(e)
	}

	entity, e := uc.repository.FindByID(idVO)
	if e != nil {
		return NewError(NotFoundDataError)
	}

	e = entity.AddTiming(
		input.DayOfWeek,
		input.HourTime,
	)
	if e != nil {
		return NewError(RegularTimingDuplicateError)
	}

	// 更新情報を保存
	e = uc.repository.Update(entity)
	if e != nil {
		return NewByError(e)
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
	id uint,
	input RegularTimingInput,
) (err AppError) {

	idVO, e := domain.NewPostMessageID(id)
	if e != nil {
		return NewByError(e)
	}

	entity, e := uc.repository.FindByID(idVO)
	if e != nil {
		return NewError(NotFoundDataError)
	}

	e = entity.RemoveTiming(
		input.DayOfWeek,
		input.HourTime,
	)
	if e != nil {
		return NewByError(e)
	}

	// 更新情報を保存
	e = uc.repository.Update(entity)
	if e != nil {
		return NewByError(e)
	}

	return err
}

// --- RegularPostSendUseCase ---

// RegularPostSendUseCase RegularPost配信ユースケース
type RegularPostSendUseCase struct {
	repository domain.RegularPostRepository
	adapter    domain.DiscordMessageAdapter
}

// NewRegularPostSendUseCase コンストラクタ
func NewRegularPostSendUseCase(
	repository domain.RegularPostRepository,
	adapter domain.DiscordMessageAdapter,
) RegularPostSendUseCase {
	return RegularPostSendUseCase{
		repository,
		adapter,
	}
}

// Execute RegularPost配信を実行
func (uc RegularPostSendUseCase) Execute(
	now time.Time,
) (err AppError) {
	timing := domain.CreateRegularTimingByTime(now)
	messages, e := uc.repository.SendList(timing)
	if e != nil {
		return NewByError(e)
	}

	for _, entity := range messages {

		// 配信状態にする
		send, e := entity.Send(now)
		if e != nil {
			return NewByError(e)
		}

		// 配信状態のメッセージをディスコードで配信
		e = uc.adapter.SendMessage(entity.Bot(), send.Message())
		if e != nil {
			return NewByError(e)
		}

		// 配信状態の情報を保存
		e = uc.repository.Update(entity)
		if e != nil {
			return NewByError(e)
		}

	}

	return err
}

// --- RegularTimingInput ---

// RegularTimingStoreInput
type RegularTimingStoreInput struct {
	DayOfWeek string
	HourTime  time.Time
}
