package application

import (
	"time"

	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	query "github.com/takemo101/dc-scheduler/pkg/application/query"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- PostMessageDeleteUseCase ---

// PostMessageDeleteUseCase PostMessage削除ユースケース
type PostMessageDeleteUseCase struct {
	repository domain.PostMessageRepository
}

// NewPostMessageDeleteUseCase コンストラクタ
func NewPostMessageDeleteUseCase(
	repository domain.PostMessageRepository,
) PostMessageDeleteUseCase {
	return PostMessageDeleteUseCase{
		repository,
	}
}

// Execute PostMessage削除を実行
func (uc PostMessageDeleteUseCase) Execute(
	context domain.UserAuthContext,
	id uint,
) (err common.AppError) {
	deleteID, e := domain.NewPostMessageID(id)
	if e != nil {
		return common.NewByError(e)
	}

	entity, e := uc.repository.FindBaseByID(deleteID)
	if e != nil {
		return common.NewByError(e)
	}

	// ポリシーチェック
	policy := domain.NewUserMessagePolicy(context)
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

	return err
}

// --- SentMessageHistoryInput ---

// SentMessageHistoryInput SentMessage一覧取得（配信メッセージ履歴取得）DTO
type SentMessageHistoryInput struct {
	Page  int
	Limit int
}

// --- SentMessageHistoryUseCase ---

// SentMessageHistoryUseCase PostMessage削除ユースケース
type SentMessageHistoryUseCase struct {
	query query.SentMessageQuery
}

// NewSentMessageHistoryUseCase コンストラクタ
func NewSentMessageHistoryUseCase(
	query query.SentMessageQuery,
) SentMessageHistoryUseCase {
	return SentMessageHistoryUseCase{
		query,
	}
}

// Execute SentMessage一覧取得（配信メッセージ履歴取得）を実行
func (uc SentMessageHistoryUseCase) Execute(
	context domain.UserAuthContext,
	input SentMessageHistoryInput,
) (paginator query.SentMessageSearchPaginatorDTO, err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return paginator, common.NewByError(e)
	}

	parameter := query.SentMessageSearchParameterDTO{
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

// --- ImmediatePostSearchInput ---

// ImmediatePostSearchInput ImmediatePost一覧取得DTO
type ImmediatePostSearchInput struct {
	Page  int
	Limit int
}

// --- ImmediatePostSearchUseCase ---

// ImmediatePostSearchUseCase ImmediatePost一覧ユースケース
type ImmediatePostSearchUseCase struct {
	query query.ImmediatePostQuery
}

// NewImmediatePostSearchUseCase コンストラクタ
func NewImmediatePostSearchUseCase(
	query query.ImmediatePostQuery,
) ImmediatePostSearchUseCase {
	return ImmediatePostSearchUseCase{
		query,
	}
}

// Execute ImmediatePost一覧取得を実行
func (uc ImmediatePostSearchUseCase) Execute(
	context domain.UserAuthContext,
	input ImmediatePostSearchInput,
) (paginator query.ImmediatePostSearchPaginatorDTO, err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return paginator, common.NewByError(e)
	}

	parameter := query.ImmediatePostSearchParameterDTO{
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

// --- ImmediatePostStoreInput ---

// ImmediatePostStoreInput ImmediatePost追加DTO
type ImmediatePostStoreInput struct {
	Message string
}

// --- ImmediatePostStoreUseCase ---

// ImmediatePostStoreUseCase ImmediatePost追加ユースケース
type ImmediatePostStoreUseCase struct {
	repository    domain.ImmediatePostRepository
	botRepository domain.BotRepository
	adapter       domain.DiscordMessageAdapter
}

// NewImmediatePostStoreUseCase コンストラクタ
func NewImmediatePostStoreUseCase(
	repository domain.ImmediatePostRepository,
	botRepository domain.BotRepository,
	adapter domain.DiscordMessageAdapter,
) ImmediatePostStoreUseCase {
	return ImmediatePostStoreUseCase{
		repository,
		botRepository,
		adapter,
	}
}

// Execute ImmediatePost追加＆即時配信を実行
func (uc ImmediatePostStoreUseCase) Execute(
	context domain.UserAuthContext,
	botID uint,
	input ImmediatePostStoreInput,
) (id uint, err common.AppError) {
	botIDVO, e := domain.NewBotID(botID)
	if e != nil {
		return id, common.NewByError(e)
	}

	bot, e := uc.botRepository.FindByID(botIDVO)
	if e != nil {
		return id, common.NewError(common.NotFoundDataError)
	}

	// ポリシーチェック
	policy := domain.NewUserMessagePolicy(context)
	ok, e := policy.Create(bot)
	if e != nil {
		return id, common.NewByError(e)
	} else if !ok {
		return id, common.NewError(common.NotTargetOwnerError)
	}

	nextID, e := uc.repository.NextIdentity()
	if e != nil {
		return id, common.NewByError(e)
	}

	entity, e := domain.CreateImmediatePost(
		nextID.Value(),
		input.Message,
		bot,
	)
	if e != nil {
		return id, common.NewByError(e)
	}

	// 配信状態にする
	send, e := entity.Send(time.Now())
	if e != nil {
		return id, common.NewByError(e)
	}

	// 配信状態のメッセージをディスコードで配信
	e = uc.adapter.SendMessage(entity.Bot(), send.Message())
	if e != nil {
		return id, common.NewByError(e)
	}

	// 配信状態の情報を保存
	storeID, e := uc.repository.Store(entity)
	if e != nil {
		return id, common.NewByError(e)
	}

	return storeID.Value(), err
}

// --- SchedulePostSearchInput ---

// SchedulePostSearchInput SchedulePost一覧取得DTO
type SchedulePostSearchInput struct {
	Page  int
	Limit int
}

// --- SchedulePostSearchUseCase ---

// SchedulePostSearchUseCase SchedulePost一覧ユースケース
type SchedulePostSearchUseCase struct {
	query query.SchedulePostQuery
}

// NewSchedulePostSearchUseCase コンストラクタ
func NewSchedulePostSearchUseCase(
	query query.SchedulePostQuery,
) SchedulePostSearchUseCase {
	return SchedulePostSearchUseCase{
		query,
	}
}

// Execute SchedulePost一覧取得を実行
func (uc SchedulePostSearchUseCase) Execute(
	context domain.UserAuthContext,
	input SchedulePostSearchInput,
) (paginator query.SchedulePostSearchPaginatorDTO, err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return paginator, common.NewByError(e)
	}

	parameter := query.SchedulePostSearchParameterDTO{
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

// --- SchedulePostStoreInput ---

// SchedulePostStoreInput SchedulePost追加DTO
type SchedulePostStoreInput struct {
	Message       string
	ReservationAt time.Time
}

// --- SchedulePostStoreUseCase ---

// SchedulePostStoreUseCase SchedulePost追加ユースケース
type SchedulePostStoreUseCase struct {
	repository    domain.SchedulePostRepository
	botRepository domain.BotRepository
}

// NewSchedulePostStoreUseCase コンストラクタ
func NewSchedulePostStoreUseCase(
	repository domain.SchedulePostRepository,
	botRepository domain.BotRepository,
) SchedulePostStoreUseCase {
	return SchedulePostStoreUseCase{
		repository,
		botRepository,
	}
}

// Execute SchedulePost追加を実行
func (uc SchedulePostStoreUseCase) Execute(
	context domain.UserAuthContext,
	botID uint,
	input SchedulePostStoreInput,
) (id uint, err common.AppError) {
	botIDVO, e := domain.NewBotID(botID)
	if e != nil {
		return id, common.NewByError(e)
	}

	bot, e := uc.botRepository.FindByID(botIDVO)
	if e != nil {
		return id, common.NewError(common.NotFoundDataError)
	}

	// ポリシーチェック
	policy := domain.NewUserMessagePolicy(context)
	ok, e := policy.Create(bot)
	if e != nil {
		return id, common.NewByError(e)
	} else if !ok {
		return id, common.NewError(common.NotTargetOwnerError)
	}

	nextID, e := uc.repository.NextIdentity()
	if e != nil {
		return id, common.NewByError(e)
	}

	entity, e := domain.CreateSchedulePost(
		nextID.Value(),
		input.Message,
		input.ReservationAt,
		bot,
		time.Now(),
	)
	if e != nil {
		return id, common.NewByError(e)
	}

	// 配信状態の情報を保存
	storeID, e := uc.repository.Store(entity)
	if e != nil {
		return id, common.NewByError(e)
	}

	return storeID.Value(), err
}

// --- SchedulePostEditFormUseCase ---

// SchedulePostEditFormUseCase SchedulePost編集フォームユースケース
type SchedulePostEditFormUseCase struct {
	repository domain.SchedulePostRepository
	query      query.SchedulePostQuery
}

// NewSchedulePostEditFormUseCase コンストラクタ
func NewSchedulePostEditFormUseCase(
	repository domain.SchedulePostRepository,
	query query.SchedulePostQuery,
) SchedulePostEditFormUseCase {
	return SchedulePostEditFormUseCase{
		repository,
		query,
	}
}

// Execute フォーム表示のためのSchedulePost取得を実行
func (uc SchedulePostEditFormUseCase) Execute(
	context domain.UserAuthContext,
	id uint,
) (detail query.SchedulePostDetailDTO, err common.AppError) {
	findID, e := domain.NewPostMessageID(id)
	if e != nil {
		return detail, common.NewByError(e)
	}

	entity, e := uc.repository.FindByID(findID)
	if e != nil {
		return detail, common.NewByError(e)
	} else if entity.IsSended() {
		return detail, common.NewError(common.NotFoundDataError)
	}

	// ポリシーチェック
	policy := domain.NewUserMessagePolicy(context)
	ok, e := policy.Detail(entity.PostMessage)
	if e != nil {
		return detail, common.NewByError(e)
	} else if !ok {
		return detail, common.NewError(common.NotTargetOwnerError)
	}

	detail, e = uc.query.FindByID(findID)
	if e != nil {
		return detail, common.NewByError(e)
	}

	return detail, err
}

// --- SchedulePostUpdateInput ---

// SchedulePostUpdateInput SchedulePost追加DTO
type SchedulePostUpdateInput struct {
	Message       string
	ReservationAt time.Time
}

// --- SchedulePostUpdateUseCase ---

// SchedulePostUpdateUseCase SchedulePost追加ユースケース
type SchedulePostUpdateUseCase struct {
	repository domain.SchedulePostRepository
}

// NewSchedulePostUpdateUseCase コンストラクタ
func NewSchedulePostUpdateUseCase(
	repository domain.SchedulePostRepository,
) SchedulePostUpdateUseCase {
	return SchedulePostUpdateUseCase{
		repository,
	}
}

// Execute SchedulePost更新を実行
func (uc SchedulePostUpdateUseCase) Execute(
	context domain.UserAuthContext,
	id uint,
	input SchedulePostUpdateInput,
) (err common.AppError) {
	idVO, e := domain.NewPostMessageID(id)
	if e != nil {
		return common.NewByError(e)
	}

	entity, e := uc.repository.FindByID(idVO)
	if e != nil {
		return common.NewError(common.NotFoundDataError)
	}

	// ポリシーチェック
	policy := domain.NewUserMessagePolicy(context)
	ok, e := policy.Update(entity.PostMessage)
	if e != nil {
		return common.NewByError(e)
	} else if !ok {
		return common.NewError(common.NotTargetOwnerError)
	}

	e = entity.Update(
		input.Message,
		input.ReservationAt,
		time.Now(),
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
