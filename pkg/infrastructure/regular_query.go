package infrastructure

import (
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/query"
	"github.com/takemo101/dc-scheduler/pkg/domain"
	"gorm.io/gorm"
)

// --- RegularPostQuery ---

// RegularPostQuery application.PostMessageQueryの実装
type RegularPostQuery struct {
	db     core.Database
	upload UploadAdapter
}

// NewRegularPostQuery コンストラクタ
func NewRegularPostQuery(
	db core.Database,
	upload UploadAdapter,
) application.RegularPostQuery {
	return RegularPostQuery{
		db,
		upload,
	}
}

// Search RegularPostの一覧取得
func (query RegularPostQuery) Search(parameter application.RegularPostSearchParameterDTO) (dto application.RegularPostSearchPaginatorDTO, err error) {
	var models []PostMessage

	paging := NewGormPaging(
		query.db.GormDB.Preload("RegularTimings").Preload("Bot").Where("message_type = ?", domain.MessageTypeRegularPost),
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderByType.ToQuery(parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	messages := make([]application.RegularPostDetailDTO, len(models))
	for i, m := range models {
		messages[i] = CreateRegularPostDetailDTOFromModel(query.upload, m)
	}

	dto.RegularPosts = messages

	return dto, err
}

// FindByID RegularPostの詳細取得
func (query RegularPostQuery) FindByID(id domain.PostMessageID) (dto application.RegularPostDetailDTO, err error) {
	model := PostMessage{}

	if err = query.db.GormDB.Where("id = ? AND message_type = ?", id.Value(), domain.MessageTypeRegularPost).Preload("RegularTimings", func(db *gorm.DB) *gorm.DB {
		return db.Order("regular_timings.hour_time ASC")
	}).Preload("Bot").First(&model).Error; err != nil {
		return dto, err
	}

	// DayOfWeek順に並び替え
	length := len(model.RegularTimings)
	if length > 0 {
		timings := make([]RegularTiming, length)
		counter := 0
		for _, week := range domain.DayOfWeeks() {
			for _, tm := range model.RegularTimings {
				if tm.DayOfWeek.Equals(week) {
					timings[counter] = tm
					counter++
				}
			}
		}
		model.RegularTimings = timings
	}

	return CreateRegularPostDetailDTOFromModel(query.upload, model), err
}

// CreateRegularTimingDTOFromModel RegularTimingからRegularTimingDTOを生成する
func CreateRegularTimingDTOFromModel(model RegularTiming) application.RegularTimingDTO {
	return application.RegularTimingDTO{
		ID:        model.ID,
		DayOfWeek: model.DayOfWeek,
		HourTime:  model.HourTime,
	}
}

// CreateRegularPostDetailDTOFromModel PostMessageからRegularPostDetailDTOを生成する
func CreateRegularPostDetailDTOFromModel(
	upload UploadAdapter,
	model PostMessage,
) application.RegularPostDetailDTO {
	bot := model.Bot

	timings := make([]application.RegularTimingDTO, len(model.RegularTimings))
	for i, tm := range model.RegularTimings {
		timings[i] = CreateRegularTimingDTOFromModel(tm)
	}

	return application.RegularPostDetailDTO{
		ID:             model.ID,
		Message:        model.Message,
		Bot:            CreateBotDetailDTOFromModel(upload, bot),
		Active:         model.Active.Bool,
		RegularTimings: timings,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}
}
