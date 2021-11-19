package application

import (
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
	context domain.UserAuthContext,
	input ApiPostSearchInput,
) (paginator query.ApiPostSearchPaginatorDTO, err common.AppError) {
	auth, e := context.UserAuth()
	if e != nil {
		return paginator, common.NewByError(e)
	}

	parameter := query.ApiPostSearchParameterDTO{
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

// --- ApiPostStoreUseCase ---

// ApiPostStoreUseCase ApiPost追加ユースケース
type ApiPostStoreUseCase struct {
	repository    domain.ApiPostRepository
	botRepository domain.BotRepository
}

// NewApiPostStoreUseCase コンストラクタ
func NewApiPostStoreUseCase(
	repository domain.ApiPostRepository,
	botRepository domain.BotRepository,
) ApiPostStoreUseCase {
	return ApiPostStoreUseCase{
		repository,
		botRepository,
	}
}

// Execute ApiPost追加＆即時配信を実行
func (uc ApiPostStoreUseCase) Execute(
	context domain.UserAuthContext,
	botID uint,
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

	entity, e := domain.CreateApiPost(
		nextID.Value(),
		bot,
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
