package infrastructure

import (
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/application"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- BotQuery ---

// BotQuery application.BotQueryの実装
type BotQuery struct {
	db     core.Database
	upload UploadAdapter
}

// NewBotQuery コンストラクタ
func NewBotQuery(
	db core.Database,
	upload UploadAdapter,
) application.BotQuery {
	return BotQuery{
		db,
		upload,
	}
}

// Search Botの一覧取得
func (query BotQuery) Search(parameter application.BotSearchParameterDTO) (dto application.BotSearchPaginatorDTO, err error) {
	var models []Bot

	paging := NewGormPaging(
		query.db.GormDB,
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderBy},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	Bots := make([]application.BotDetailDTO, len(models))
	for i, m := range models {
		Bots[i] = CreateBotDetailDTOFromModel(query.upload, m)
	}

	dto.Bots = Bots

	return dto, err
}

// FindByID Botの詳細取得
func (query BotQuery) FindByID(id domain.BotID) (dto application.BotDetailDTO, err error) {
	model := Bot{}

	if err = query.db.GormDB.Where("id = ?", id.Value()).First(&model).Error; err != nil {
		return dto, err
	}

	return CreateBotDetailDTOFromModel(query.upload, model), err
}

// CreateBotDetailDTOFromModel BotからBotDetailDTOを生成する
func CreateBotDetailDTOFromModel(
	upload UploadAdapter,
	model Bot,
) application.BotDetailDTO {
	var avatorURL, avatorPath string
	if model.Avator != "" {
		avatorURL = upload.ToURL(model.Avator)
		avatorPath = model.Avator
	}

	return application.BotDetailDTO{
		ID:         model.ID,
		Name:       model.Name,
		AvatorURL:  avatorURL,
		AvatorPath: avatorPath,
		Webhook:    model.Webhook,
		Active:     model.Active,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}
