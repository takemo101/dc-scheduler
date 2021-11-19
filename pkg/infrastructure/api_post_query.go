package infrastructure

import (
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/query"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- ApiPostQuery ---

// ApiPostQuery application.ApiPostQueryの実装
type ApiPostQuery struct {
	db     core.Database
	upload UploadAdapter
}

// NewApiPostQuery コンストラクタ
func NewApiPostQuery(
	db core.Database,
	upload UploadAdapter,
) application.ApiPostQuery {
	return ApiPostQuery{
		db,
		upload,
	}
}

// Search ApiPostの一覧取得
func (query ApiPostQuery) Search(parameter application.ApiPostSearchParameterDTO) (dto application.ApiPostSearchPaginatorDTO, err error) {
	var models []PostMessage

	paging := NewGormPaging(
		query.db.GormDB.Preload("Bot.User").Preload("ApiKey").Where("message_type = ?", domain.MessageTypeApiPost),
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderByType.ToQuery(parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	messages := make([]application.ApiPostDetailDTO, len(models))
	for i, m := range models {
		messages[i] = CreateApiPostDetailDTOFromModel(query.upload, m)
	}

	dto.ApiPosts = messages

	return dto, err
}

// SearchByUserID UserのApiPost一覧取得
func (query ApiPostQuery) SearchByUserID(parameter application.ApiPostSearchParameterDTO, userID domain.UserID) (dto application.ApiPostSearchPaginatorDTO, err error) {
	var models []PostMessage

	paging := NewGormPaging(
		query.db.GormDB.Preload("Bot").Preload("ApiKey").Joins(
			"JOIN bots ON bots.id = post_messages.bot_id AND bots.user_id = ?",
			userID.Value(),
		).Where(
			"message_type = ?",
			domain.MessageTypeApiPost,
		),
		parameter.Page,
		parameter.Limit,
		[]string{"post_messages." + parameter.OrderByType.ToQuery(parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	messages := make([]application.ApiPostDetailDTO, len(models))
	for i, m := range models {
		messages[i] = CreateApiPostDetailDTOFromModel(query.upload, m)
	}

	dto.ApiPosts = messages

	return dto, err
}

// CreateApiPostDetailDTOFromModel ApiPostからApiPostDetailDTOを生成する
func CreateApiPostDetailDTOFromModel(
	upload UploadAdapter,
	model PostMessage,
) application.ApiPostDetailDTO {
	bot := model.Bot

	return application.ApiPostDetailDTO{
		ID:        model.ID,
		Message:   model.Message,
		ApiKey:    model.ApiKey.Key,
		Bot:       CreateBotDetailDTOFromModel(upload, bot),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
