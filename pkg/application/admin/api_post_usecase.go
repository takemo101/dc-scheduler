package application

import (
	"time"

	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	query "github.com/takemo101/dc-scheduler/pkg/application/query"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- ApiPostSearchInput ---

// ApiPostSearchInput ApiPost一覧取得DTO
type ApiPostSearchInput struct {
	Page  int
	Limit int
}

// --- ApiPostSearchUseCase ---

// ApiPostSearchUseCase ApiPost一覧ユースケース
type ApiPostSearchUseCase struct {
	query query.ApiPostQuery
}

// NewApiPostSearchUseCase コンストラクタ
func NewApiPostSearchUseCase(
	query query.ApiPostQuery,
) ApiPostSearchUseCase {
	return ApiPostSearchUseCase{
		query,
	}
}

// Execute ApiPost一覧取得を実行
func (uc ApiPostSearchUseCase) Execute(
	input ApiPostSearchInput,
) (paginator query.ApiPostSearchPaginatorDTO, err common.AppError) {

	parameter := query.ApiPostSearchParameterDTO{
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

// --- ApiPostSendUseCase ---

// ApiPostSendUseCase ApiPost配信ユースケース
type ApiPostSendUseCase struct {
	repository domain.ApiPostRepository
	adapter    domain.DiscordMessageAdapter
}

// NewApiPostSendUseCase コンストラクタ
func NewApiPostSendUseCase(
	repository domain.ApiPostRepository,
	adapter domain.DiscordMessageAdapter,
) ApiPostSendUseCase {
	return ApiPostSendUseCase{
		repository,
		adapter,
	}
}

// Execute ApiPost配信を実行
func (uc ApiPostSendUseCase) Execute(
	key string,
	message string,
	now time.Time,
) (err common.AppError) {

	apiKey, e := domain.NewMessageApiKey(key)
	if e != nil {
		return common.NewByError(e)
	}

	entity, e := uc.repository.FindByApiKey(apiKey)
	if e != nil {
		return common.NewError(common.NotFoundDataError)
	}

	// 配信状態にする
	send, e := entity.Send(message, now)
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

	return err
}
