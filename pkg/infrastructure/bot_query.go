package infrastructure

import (
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/query"
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
		[]string{parameter.OrderByType.ToQuery(parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	bots := make([]application.BotDetailDTO, len(models))
	for i, m := range models {
		bots[i] = CreateBotDetailDTOFromModel(query.upload, m)
	}

	dto.Bots = bots

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
	var avatarURL, avatarPath string
	if model.Atatar != "" {
		avatarURL = upload.ToURL(model.Atatar)
		avatarPath = model.Atatar
	}

	return application.BotDetailDTO{
		ID:         model.ID,
		Name:       model.Name,
		AtatarURL:  avatarURL,
		AtatarPath: avatarPath,
		Webhook:    model.Webhook,
		Active:     model.Active.Bool,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}
