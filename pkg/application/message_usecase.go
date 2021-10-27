package application

import (
	"time"

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
	query PostMessageQuery
}

// NewPostMessageSearchUseCase コンストラクタ
func NewPostMessageSearchUseCase(
	query PostMessageQuery,
) PostMessageSearchUseCase {
	return PostMessageSearchUseCase{
		query,
	}
}

// Execute PostMessage一覧取得を実行
func (uc PostMessageSearchUseCase) Execute(
	input PostMessageSearchInput,
) (paginator PostMessageSearchPaginatorDTO, err AppError) {

	parameter := PostMessageSearchParameterDTO{
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

// --- PostMessageCreateFormUseCase ---

// PostMessageCreateFormUseCase Bot作成フォームユースケース
type PostMessageCreateFormUseCase struct {
	repository domain.BotRepository
	query      BotQuery
}

// NewPostMessageCreateFormUseCase コンストラクタ
func NewPostMessageCreateFormUseCase(
	repository domain.BotRepository,
	query BotQuery,
) PostMessageCreateFormUseCase {
	return PostMessageCreateFormUseCase{
		repository,
		query,
	}
}

// Execute フォーム表示のためのBot取得を実行
func (uc PostMessageCreateFormUseCase) Execute(botID uint) (detail BotDetailDTO, err AppError) {
	findID, e := domain.NewBotID(botID)
	if e != nil {
		return detail, NewError(NotFoundDataError)
	}

	detail, e = uc.query.FindByID(findID)
	if e != nil {
		return detail, NewByError(e)
	}

	return detail, err
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
func (uc PostMessageDeleteUseCase) Execute(id uint) (err AppError) {

	deleteID, e := domain.NewPostMessageID(id)
	if e != nil {
		return NewByError(e)
	}

	e = uc.repository.Delete(deleteID)
	if e != nil {
		return NewByError(e)
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
	query SentMessageQuery
}

// NewSentMessageHistoryUseCase コンストラクタ
func NewSentMessageHistoryUseCase(
	query SentMessageQuery,
) SentMessageHistoryUseCase {
	return SentMessageHistoryUseCase{
		query,
	}
}

// Execute SentMessage一覧取得（配信メッセージ履歴取得）を実行
func (uc SentMessageHistoryUseCase) Execute(
	input SentMessageHistoryInput,
) (paginator SentMessageSearchPaginatorDTO, err AppError) {

	parameter := SentMessageSearchParameterDTO{
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

// --- ImmediatePostSearchInput ---

// ImmediatePostSearchInput ImmediatePost一覧取得DTO
type ImmediatePostSearchInput struct {
	Page  int
	Limit int
}

// --- ImmediatePostSearchUseCase ---

// ImmediatePostSearchUseCase ImmediatePost一覧ユースケース
type ImmediatePostSearchUseCase struct {
	query ImmediatePostQuery
}

// NewImmediatePostSearchUseCase コンストラクタ
func NewImmediatePostSearchUseCase(
	query ImmediatePostQuery,
) ImmediatePostSearchUseCase {
	return ImmediatePostSearchUseCase{
		query,
	}
}

// Execute ImmediatePost一覧取得を実行
func (uc ImmediatePostSearchUseCase) Execute(
	input ImmediatePostSearchInput,
) (paginator ImmediatePostSearchPaginatorDTO, err AppError) {

	parameter := ImmediatePostSearchParameterDTO{
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
	botID uint,
	input ImmediatePostStoreInput,
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

	entity, e := domain.CreateImmediatePost(
		nextID.Value(),
		input.Message,
		bot,
	)
	if e != nil {
		return id, NewByError(e)
	}

	// 配信状態にする
	send, e := entity.Send(time.Now())
	if e != nil {
		return id, NewByError(e)
	}

	// 配信状態のメッセージをディスコードで配信
	e = uc.adapter.SendMessage(entity.Bot(), send.Message())
	if e != nil {
		return id, NewByError(e)
	}

	// 配信状態の情報を保存
	storeID, e := uc.repository.Store(entity)
	if e != nil {
		return id, NewByError(e)
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
	query SchedulePostQuery
}

// NewSchedulePostSearchUseCase コンストラクタ
func NewSchedulePostSearchUseCase(
	query SchedulePostQuery,
) SchedulePostSearchUseCase {
	return SchedulePostSearchUseCase{
		query,
	}
}

// Execute SchedulePost一覧取得を実行
func (uc SchedulePostSearchUseCase) Execute(
	input SchedulePostSearchInput,
) (paginator SchedulePostSearchPaginatorDTO, err AppError) {

	parameter := SchedulePostSearchParameterDTO{
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

	entity, e := domain.CreateSchedulePost(
		nextID.Value(),
		input.Message,
		input.ReservationAt,
		bot,
		time.Now(),
	)
	if e != nil {
		return id, NewByError(e)
	}

	// 配信状態の情報を保存
	storeID, e := uc.repository.Store(entity)
	if e != nil {
		return id, NewByError(e)
	}

	return storeID.Value(), err
}

// --- SchedulePostEditFormUseCase ---

// SchedulePostEditFormUseCase SchedulePost編集フォームユースケース
type SchedulePostEditFormUseCase struct {
	repository domain.SchedulePostRepository
	query      SchedulePostQuery
}

// NewPostMessageCreateFormUseCase コンストラクタ
func NewSchedulePostEditFormUseCase(
	repository domain.SchedulePostRepository,
	query SchedulePostQuery,
) SchedulePostEditFormUseCase {
	return SchedulePostEditFormUseCase{
		repository,
		query,
	}
}

// Execute フォーム表示のためのSchedulePost取得を実行
func (uc SchedulePostEditFormUseCase) Execute(id uint) (detail SchedulePostDetailDTO, err AppError) {
	findID, e := domain.NewPostMessageID(id)
	if e != nil {
		return detail, NewByError(e)
	}

	entity, e := uc.repository.FindByID(findID)
	if e != nil {
		return detail, NewByError(e)
	} else if entity.IsSended() {
		return detail, NewError(NotFoundDataError)
	}

	detail, e = uc.query.FindByID(findID)
	if e != nil {
		return detail, NewByError(e)
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
		input.ReservationAt,
		time.Now(),
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
) (err AppError) {
	messages, e := uc.repository.SendList(domain.NewMessageSendedAt(now))
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
