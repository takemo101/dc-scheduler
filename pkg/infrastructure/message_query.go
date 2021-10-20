package infrastructure

import (
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/application"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- PostMessageQuery ---

// PostMessageQuery application.PostMessageQueryの実装
type PostMessageQuery struct {
	db     core.Database
	upload UploadAdapter
}

// NewPostMessageQuery コンストラクタ
func NewPostMessageQuery(
	db core.Database,
	upload UploadAdapter,
) application.PostMessageQuery {
	return PostMessageQuery{
		db,
		upload,
	}
}

// Search PostMessageの一覧取得
func (query PostMessageQuery) Search(parameter application.PostMessageSearchParameterDTO) (dto application.PostMessageSearchPaginatorDTO, err error) {
	var models []PostMessage

	paging := NewGormPaging(
		query.db.GormDB.Preload("Bot"),
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderBy},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	PostMessages := make([]application.PostMessageDetailDTO, len(models))
	for i, m := range models {
		PostMessages[i] = CreatePostMessageDetailDTOFromModel(query.upload, m)
	}

	dto.PostMessages = PostMessages

	return dto, err
}

// FindByID PostMessageの詳細取得
func (query PostMessageQuery) FindByID(id domain.PostMessageID) (dto application.PostMessageDetailDTO, err error) {
	model := PostMessage{}

	if err = query.db.GormDB.Where("id = ?", id.Value()).Preload("Bot").First(&model).Error; err != nil {
		return dto, err
	}

	return CreatePostMessageDetailDTOFromModel(query.upload, model), err
}

// CreatePostMessageDetailDTOFromModel PostMessageからPostMessageDetailDTOを生成する
func CreatePostMessageDetailDTOFromModel(
	upload UploadAdapter,
	model PostMessage,
) application.PostMessageDetailDTO {
	bot := model.Bot

	return application.PostMessageDetailDTO{
		ID:          model.ID,
		Message:     model.Message,
		MessageType: model.MessageType,
		Bot:         CreateBotDetailDTOFromModel(upload, bot),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
