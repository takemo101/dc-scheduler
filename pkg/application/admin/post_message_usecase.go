package application

import (
	"time"

	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	query "github.com/takemo101/dc-scheduler/pkg/application/query"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- PostMessageSearchInput ---

// PostMessageSearchInput PostMessage一覧取得DTO
type PostMessageSearchInput struct {
	Page  int
	Limit int
}

// --- PostMessageSearchUseCase ---

// PostMessageSearchUseCase PostMessage一覧ユースケース
type PostMessageSearchUseCase struct {
	query query.PostMessageQuery
}

// NewPostMessageSearchUseCase コンストラクタ
func NewPostMessageSearchUseCase(
	query query.PostMessageQuery,
) PostMessageSearchUseCase {
	return PostMessageSearchUseCase{
		query,
	}
}

// Execute PostMessage一覧取得を実行
func (uc PostMessageSearchUseCase) Execute(
	input PostMessageSearchInput,
) (paginator query.PostMessageSearchPaginatorDTO, err common.AppError) {

	parameter := query.PostMessageSearchParameterDTO{
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
func (uc PostMessageDeleteUseCase) Execute(id uint) (err common.AppError) {

	deleteID, e := domain.NewPostMessageID(id)
	if e != nil {
		return common.NewByError(e)
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
	input SentMessageHistoryInput,
) (paginator query.SentMessageSearchPaginatorDTO, err common.AppError) {

	parameter := query.SentMessageSearchParameterDTO{
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
	input ImmediatePostSearchInput,
) (paginator query.ImmediatePostSearchPaginatorDTO, err common.AppError) {

	parameter := query.ImmediatePostSearchParameterDTO{
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
	input SchedulePostSearchInput,
) (paginator query.SchedulePostSearchPaginatorDTO, err common.AppError) {

	parameter := query.SchedulePostSearchParameterDTO{
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

// NewPostMessageCreateFormUseCase コンストラクタ
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
func (uc SchedulePostEditFormUseCase) Execute(id uint) (detail query.SchedulePostDetailDTO, err common.AppError) {
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

// --- SchedulePostSendUseCase ---

// SchedulePostSendUseCase SchedulePost配信ユースケース
type SchedulePostSendUseCase struct {
	repository domain.SchedulePostRepository
	adapter    domain.DiscordMessageAdapter
}

// NewSchedulePostSendUseCase コンストラクタ
func NewSchedulePostSendUseCase(
	repository domain.SchedulePostRepository,
	adapter domain.DiscordMessageAdapter,
) SchedulePostSendUseCase {
	return SchedulePostSendUseCase{
		repository,
		adapter,
	}
}

// Execute SchedulePost配信を実行
func (uc SchedulePostSendUseCase) Execute(
	now time.Time,
) (err common.AppError) {
	messages, e := uc.repository.SendList(domain.NewJustMessageSendedAt(now))
	if e != nil {
		return common.NewByError(e)
	}

	for _, entity := range messages {

		// 配信状態にする
		send, e := entity.Send(now)
		if e != nil {
			return common.NewByError(e)
		}

		// 配信状態のメッセージをディスコードで配信
		e = uc.adapter.SendMessage(entity.Bot(), send.Message())
		if e != nil {
			return common.NewByError(e)
		}

		// 配信状態の情報を保存
		e = uc.repository.Update(entity)
		if e != nil {
			return common.NewByError(e)
		}

	}

	return err
}
