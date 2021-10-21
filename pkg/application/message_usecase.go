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

	entity, e := uc.repository.FindByID(findID)
	if e != nil {
		return detail, NewByError(e)
	} else if !entity.IsActive() {
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

// SentMessageHistoryInput SentMessage一覧取得（送信メッセージ履歴取得）DTO
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

// Execute SentMessage一覧取得（送信メッセージ履歴取得）を実行
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

// Execute ImmediatePost追加＆即時送信を実行
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

	// 送信状態にする
	send, e := entity.Send(time.Now())
	if e != nil {
		return id, NewByError(e)
	}

	// 送信状態のメッセージをディスコードで送信
	e = uc.adapter.SendMessage(entity.Bot(), send.Message())
	if e != nil {
		return id, NewByError(e)
	}

	// 送信状態の情報を保存
	storeID, e := uc.repository.Store(entity)
	if e != nil {
		return id, NewByError(e)
	}

	return storeID.Value(), err
}
