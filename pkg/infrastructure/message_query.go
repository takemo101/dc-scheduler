package infrastructure

import (
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/query"
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
		[]string{parameter.OrderByType.ToQuery(parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	messages := make([]application.PostMessageDetailDTO, len(models))
	for i, m := range models {
		messages[i] = CreatePostMessageDetailDTOFromModel(query.upload, m)
	}

	dto.PostMessages = messages

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

// --- SentMessageQuery ---

// SentMessageQuery application.SentMessageQueryの実装
type SentMessageQuery struct {
	db     core.Database
	upload UploadAdapter
}

// NewSentMessageQuery コンストラクタ
func NewSentMessageQuery(
	db core.Database,
	upload UploadAdapter,
) application.SentMessageQuery {
	return SentMessageQuery{
		db,
		upload,
	}
}

// Search SentMessageの一覧取得
func (query SentMessageQuery) Search(parameter application.SentMessageSearchParameterDTO) (dto application.SentMessageSearchPaginatorDTO, err error) {
	var models []SentMessage

	paging := NewGormPaging(
		query.db.GormDB.Preload("PostMessage.Bot"),
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderByType.ToQuery(parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	SentMessages := make([]application.SentMessageDetailDTO, len(models))
	for i, m := range models {
		SentMessages[i] = CreateSentMessageDetailDTOFromModel(query.upload, m)
	}

	dto.SentMessages = SentMessages

	return dto, err
}

// Search UserのSentMessage一覧取得
func (query SentMessageQuery) SearchByUserID(parameter application.SentMessageSearchParameterDTO, userID domain.UserID) (dto application.SentMessageSearchPaginatorDTO, err error) {
	var models []SentMessage

	paging := NewGormPaging(
		query.db.GormDB.Preload("PostMessage.Bot").Joins(
			"PostMessage.Bot",
			query.db.GormDB.Where(&Bot{UserID: userID.Value()}),
		),
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderByType.ToQuery(parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	SentMessages := make([]application.SentMessageDetailDTO, len(models))
	for i, m := range models {
		SentMessages[i] = CreateSentMessageDetailDTOFromModel(query.upload, m)
	}

	dto.SentMessages = SentMessages

	return dto, err
}

// Search SentMessageのリスト取得
func (query SentMessageQuery) RecentlyList(limit uint) (list []application.SentMessageDetailDTO, err error) {
	var models []SentMessage

	err = query.db.GormDB.Preload("PostMessage.Bot").Order("id DESC").Limit(int(limit)).Find(&models).Error
	if err != nil {
		return list, err
	}

	list = make([]application.SentMessageDetailDTO, len(models))
	for i, m := range models {
		list[i] = CreateSentMessageDetailDTOFromModel(query.upload, m)
	}

	return list, err
}

// Search UserのSentMessageリスト取得
func (query SentMessageQuery) RecentlyListByUserID(userID domain.UserID, limit uint) (list []application.SentMessageDetailDTO, err error) {
	var models []SentMessage

	err = query.db.GormDB.Preload("PostMessage.Bot").Joins(
		"PostMessage",
		query.db.GormDB.Joins("Bot", query.db.GormDB.Where(&Bot{UserID: userID.Value()})),
	).Order("sent_messages.id DESC").Limit(int(limit)).Find(&models).Error
	if err != nil {
		return list, err
	}

	list = make([]application.SentMessageDetailDTO, len(models))
	for i, m := range models {
		list[i] = CreateSentMessageDetailDTOFromModel(query.upload, m)
	}

	return list, err
}

// CreateSentMessageDetailDTOFromModel SentMessageからSentMessageDetailDTOを生成する
func CreateSentMessageDetailDTOFromModel(
	upload UploadAdapter,
	model SentMessage,
) application.SentMessageDetailDTO {
	return application.SentMessageDetailDTO{
		ID:          model.ID,
		Message:     model.Message,
		PostMessage: CreatePostMessageDetailDTOFromModel(upload, model.PostMessage),
		SendedAt:    model.SendedAt,
	}
}

// --- ImmediatePostQuery ---

// ImmediatePostQuery application.PostMessageQueryの実装
type ImmediatePostQuery struct {
	db     core.Database
	upload UploadAdapter
}

// NewImmediatePostQuery コンストラクタ
func NewImmediatePostQuery(
	db core.Database,
	upload UploadAdapter,
) application.ImmediatePostQuery {
	return ImmediatePostQuery{
		db,
		upload,
	}
}

// Search ImmediatePostの一覧取得
func (query ImmediatePostQuery) Search(parameter application.ImmediatePostSearchParameterDTO) (dto application.ImmediatePostSearchPaginatorDTO, err error) {
	var models []PostMessage

	paging := NewGormPaging(
		query.db.GormDB.Preload("Bot").Where("message_type = ?", domain.MessageTypeImmediatePost),
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderByType.ToQuery(parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	messages := make([]application.ImmediatePostDetailDTO, len(models))
	for i, m := range models {
		messages[i] = CreateImmediatePostDetailDTOFromModel(query.upload, m)
	}

	dto.ImmediatePosts = messages

	return dto, err
}

// SearchByUserID UserのImmediatePost一覧取得
func (query ImmediatePostQuery) SearchByUserID(parameter application.ImmediatePostSearchParameterDTO, userID domain.UserID) (dto application.ImmediatePostSearchPaginatorDTO, err error) {
	var models []PostMessage

	paging := NewGormPaging(
		query.db.GormDB.Preload("Bot").Joins(
			"Bot",
			query.db.GormDB.Where(&Bot{UserID: userID.Value()}),
		).Where("message_type = ?", domain.MessageTypeImmediatePost),
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderByType.ToQuery(parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	messages := make([]application.ImmediatePostDetailDTO, len(models))
	for i, m := range models {
		messages[i] = CreateImmediatePostDetailDTOFromModel(query.upload, m)
	}

	dto.ImmediatePosts = messages

	return dto, err
}

// CreateImmediatePostDetailDTOFromModel PostMessageからImmediatePostDetailDTOを生成する
func CreateImmediatePostDetailDTOFromModel(
	upload UploadAdapter,
	model PostMessage,
) application.ImmediatePostDetailDTO {
	bot := model.Bot

	return application.ImmediatePostDetailDTO{
		ID:        model.ID,
		Message:   model.Message,
		Bot:       CreateBotDetailDTOFromModel(upload, bot),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// --- SchedulePostQuery ---

// SchedulePostQuery application.PostMessageQueryの実装
type SchedulePostQuery struct {
	db     core.Database
	upload UploadAdapter
}

// NewSchedulePostQuery コンストラクタ
func NewSchedulePostQuery(
	db core.Database,
	upload UploadAdapter,
) application.SchedulePostQuery {
	return SchedulePostQuery{
		db,
		upload,
	}
}

// Search SchedulePostの一覧取得
func (query SchedulePostQuery) Search(parameter application.SchedulePostSearchParameterDTO) (dto application.SchedulePostSearchPaginatorDTO, err error) {
	var models []PostMessage

	paging := NewGormPaging(
		query.db.GormDB.Preload("ScheduleTiming").Preload("Bot.User").Where("message_type = ?", domain.MessageTypeSchedulePost).Joins("ScheduleTiming"),
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderByType.ToQuery("post_messages." + parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	messages := make([]application.SchedulePostDetailDTO, len(models))
	for i, m := range models {
		messages[i] = CreateSchedulePostDetailDTOFromModel(query.upload, m)
	}

	dto.SchedulePosts = messages

	return dto, err
}

// SearchByUserID SchedulePostの一覧取得
func (query SchedulePostQuery) SearchByUserID(parameter application.SchedulePostSearchParameterDTO, userID domain.UserID) (dto application.SchedulePostSearchPaginatorDTO, err error) {
	var models []PostMessage

	paging := NewGormPaging(
		query.db.GormDB.Preload("ScheduleTiming").Preload("Bot").Joins(
			"Bot",
			query.db.GormDB.Where(&Bot{UserID: userID.Value()}),
		).Where("message_type = ?", domain.MessageTypeSchedulePost).Joins("ScheduleTiming"),
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderByType.ToQuery("post_messages." + parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	messages := make([]application.SchedulePostDetailDTO, len(models))
	for i, m := range models {
		messages[i] = CreateSchedulePostDetailDTOFromModel(query.upload, m)
	}

	dto.SchedulePosts = messages

	return dto, err
}

// FindByID SchedulePostの詳細取得
func (query SchedulePostQuery) FindByID(id domain.PostMessageID) (dto application.SchedulePostDetailDTO, err error) {
	model := PostMessage{}

	if err = query.db.GormDB.Where("id = ? AND message_type = ?", id.Value(), domain.MessageTypeSchedulePost).Preload("ScheduleTiming").Preload("Bot.User").First(&model).Error; err != nil {
		return dto, err
	}

	return CreateSchedulePostDetailDTOFromModel(query.upload, model), err
}

// CreateSchedulePostDetailDTOFromModel PostMessageからSchedulePostDetailDTOを生成する
func CreateSchedulePostDetailDTOFromModel(
	upload UploadAdapter,
	model PostMessage,
) application.SchedulePostDetailDTO {
	bot := model.Bot

	return application.SchedulePostDetailDTO{
		ID:            model.ID,
		Message:       model.Message,
		Bot:           CreateBotDetailDTOFromModel(upload, bot),
		IsSended:      model.Sended.Bool,
		ReservationAt: model.ScheduleTiming.ReservationAt,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}
