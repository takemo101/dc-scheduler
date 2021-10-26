package infrastructure

import (
	"database/sql"
	"time"

	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/domain"
	"gorm.io/gorm"
)

// --- RegularPostRepository ---

// RegularPostRepository 予約配信Entityの永続化
type RegularPostRepository struct {
	PostMessageRepository
	db core.Database
}

// NewRegularPostRepository コンストラクタ
func NewRegularPostRepository(
	db core.Database,
	config core.Config,
) domain.RegularPostRepository {
	return RegularPostRepository{
		PostMessageRepository{
			db,
			config,
		},
		db,
	}
}

// SendList 配信可能なリストを取得
func (repo RegularPostRepository) SendList(at domain.RegularTiming) ([]domain.RegularPost, error) {
	models := []PostMessage{}

	if err := repo.db.GormDB.Preload("RegularTimings").Where("message_type = ? AND hour_time = ? AND day_of_week = ? AND post_messages.active = ? AND Bot.active = ?", domain.MessageTypeRegularPost, at.HourTime().Value(), at.DayOfWeek().Value(), true, true).Joins("Bot").Joins("left join regular_timings on regular_timings.post_message_id = post_messages.id").Find(&models).Error; err != nil {
		return []domain.RegularPost{}, err
	}

	list := make([]domain.RegularPost, len(models))

	for i, model := range models {
		list[i] = CreateRegularPostEntityFromModel(model)
	}

	return list, nil
}

// Store メッセージの追加
func (repo RegularPostRepository) Store(entity domain.RegularPost) (vo domain.PostMessageID, err error) {
	model := PostMessage{}

	model.Message = entity.Message().Value()
	model.MessageType = entity.MessageType()
	model.BotID = entity.Bot().ID().Value()
	model.Sended = sql.NullBool{Bool: false, Valid: true}
	model.Active = sql.NullBool{Bool: entity.IsActive(), Valid: true}

	// 一旦モデルを保存
	if err = repo.db.GormDB.Create(&model).Error; err != nil {
		return vo, err
	}

	// idを生成する
	return domain.NewPostMessageID(model.ID)
}

// Update メッセージの更新＆配信済みメッセージがあればそれも保存（予約配信なので）
func (repo RegularPostRepository) Update(entity domain.RegularPost) error {
	model := PostMessage{}

	model.ID = entity.ID().Value()
	model.Message = entity.Message().Value()
	model.Active = sql.NullBool{Bool: entity.IsActive(), Valid: true}

	timings := entity.Timings()
	timingModels := make([]RegularTiming, len(timings))

	for i, tm := range timings {
		timing := RegularTiming{
			PostMessageID: entity.ID().Value(),
			DayOfWeek:     tm.DayOfWeek(),
			HourTime:      tm.HourTime().Value(),
		}
		timingModels[i] = timing
	}
	model.RegularTimings = timingModels

	err := repo.db.GormDB.Transaction(func(tx *gorm.DB) (err error) {
		err = repo.db.GormDB.Where("post_message_id = ?", entity.ID().Value()).Delete(&RegularTiming{}).Error

		err = repo.db.GormDB.Updates(&model).Error
		if err != nil {
			return err
		}

		// 配信済みメッセージがあれば保存
		if entity.HasSentMessages() {
			for _, sent := range entity.SentMessages() {
				err = repo.SaveSendedMessage(entity.ID(), sent)
				if err != nil {
					return err
				}
			}
		}

		return err
	})

	return err
}

// FindByID PostMessageIDからRegularPostを取得する
func (repo RegularPostRepository) FindByID(id domain.PostMessageID) (entity domain.RegularPost, err error) {
	model := PostMessage{}

	if err = repo.db.GormDB.Where("id = ? AND message_type = ?", id.Value(), domain.MessageTypeRegularPost).Preload("RegularTimings").Preload("Bot").First(&model).Error; err != nil {
		return entity, err
	}

	return CreateRegularPostEntityFromModel(model), err
}

// CreateRegularPostEntityFromModel PostMessageからEntityを生成する
func CreateRegularPostEntityFromModel(model PostMessage) domain.RegularPost {
	timings := make([]domain.RegularTiming, len(model.RegularTimings))
	for i, tm := range model.RegularTimings {
		timings[i] = domain.NewRegularTiming(
			tm.DayOfWeek,
			tm.HourTime,
		)
	}

	return domain.NewRegularPost(
		model.ID,
		model.Message,
		CreateBotEntityFromModel(model.Bot),
		[]domain.SentMessage{},
		model.Active.Bool,
		timings,
	)
}

// --- RegularTiming ---

// RegularTiming Gormモデル
type RegularTiming struct {
	ID            uint             `gorm:"primarykey"`
	PostMessageID uint             `gorm:"index;not null"`
	DayOfWeek     domain.DayOfWeek `gorm:"type:varchar(30);index"`
	HourTime      time.Time        `gorm:"type:time;index"`
}
