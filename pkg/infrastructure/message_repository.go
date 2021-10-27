package infrastructure

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/domain"
	"gorm.io/gorm"
)

// --- PostMessageRepository ---

// PostMessageRepository 配信Entityの永続化ベース
type PostMessageRepository struct {
	db     core.Database
	config core.Config
}

// NewPostMessageRepository コンストラクタ
func NewPostMessageRepository(
	db core.Database,
	config core.Config,
) domain.PostMessageRepository {
	return PostMessageRepository{
		db,
		config,
	}
}

// SaveSendedMessage PostMessageIDからSentMessageを追加する
func (repo PostMessageRepository) SaveSendedMessage(
	id domain.PostMessageID,
	entity domain.SentMessage,
) (err error) {
	model := SentMessage{}

	model.PostMessageID = id.Value()
	model.Message = entity.Message().Value()
	model.SendedAt = entity.SendedAt().Value()

	err = repo.db.GormDB.Create(&model).Error

	return err
}

// Delete PostMessageIDからPostMessageを削除する
func (repo PostMessageRepository) Delete(id domain.PostMessageID) error {

	return repo.db.GormDB.Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Where("post_message_id = ?", id.Value()).Delete(&RegularTiming{}).Error
		if err != nil {
			return err
		}

		err = tx.Where("post_message_id = ?", id.Value()).Delete(&ScheduleTiming{}).Error
		if err != nil {
			return err
		}

		return tx.Where("id = ?", id.Value()).Delete(&PostMessage{}).Error
	})
}

// NextIdentity 次のIDを取得する
func (repo PostMessageRepository) NextIdentity() (domain.PostMessageID, error) {
	var max uint

	sql := GetNextIdentitySelectSQL(repo.config.DB.Type)
	repo.db.GormDB.Model(&PostMessage{}).Select(sql).Scan(&max)

	return domain.NewPostMessageID(max + 1)
}

// --- ImmediatePostRepository ---

// ImmediatePostRepository 即時配信Entityの永続化
type ImmediatePostRepository struct {
	PostMessageRepository
	db core.Database
}

// NewImmediatePostRepository コンストラクタ
func NewImmediatePostRepository(
	db core.Database,
	config core.Config,
) domain.ImmediatePostRepository {
	return ImmediatePostRepository{
		PostMessageRepository{
			db,
			config,
		},
		db,
	}
}

// Store メッセージの追加＆配信済みメッセージがあればそれも保存（即時配信なので）
func (repo ImmediatePostRepository) Store(entity domain.ImmediatePost) (vo domain.PostMessageID, err error) {
	model := PostMessage{}

	model.Message = entity.Message().Value()
	model.MessageType = entity.MessageType()
	model.BotID = entity.Bot().ID().Value()
	model.Sended = sql.NullBool{Bool: entity.IsSended(), Valid: true}
	model.Active = sql.NullBool{Bool: entity.IsSended(), Valid: true}

	err = repo.db.GormDB.Transaction(func(tx *gorm.DB) (err error) {
		// 一旦モデルを保存
		if err = repo.db.GormDB.Create(&model).Error; err != nil {
			return err
		}

		// idを生成する
		vo, err = domain.NewPostMessageID(model.ID)
		if err != nil {
			return err
		}

		// 配信済みメッセージがあれば保存
		if entity.HasSentMessages() {
			for _, sent := range entity.SentMessages() {
				err = repo.SaveSendedMessage(vo, sent)
				if err != nil {
					return err
				}
			}
		}

		return err
	})
	if err != nil {
		return vo, err
	}

	return vo, err
}

// FindByID PostMessageIDからImmediatePostを取得する
func (repo ImmediatePostRepository) FindByID(id domain.PostMessageID) (entity domain.ImmediatePost, err error) {
	model := PostMessage{}

	if err = repo.db.GormDB.Where("id = ? AND message_type = ?", id.Value(), domain.MessageTypeImmediatePost).Preload("Bot").First(&model).Error; err != nil {
		return entity, err
	}

	return CreateImmediatePostEntityFromModel(model), err
}

// CreateImmediatePostEntityFromModel PostMessageからEntityを生成する
func CreateImmediatePostEntityFromModel(model PostMessage) domain.ImmediatePost {
	return domain.NewImmediatePost(
		model.ID,
		model.Message,
		CreateBotEntityFromModel(model.Bot),
		[]domain.SentMessage{},
		model.Sended.Bool,
	)
}

// --- SchedulePostRepository ---

// SchedulePostRepository 予約配信Entityの永続化
type SchedulePostRepository struct {
	PostMessageRepository
	db core.Database
}

// NewSchedulePostRepository コンストラクタ
func NewSchedulePostRepository(
	db core.Database,
	config core.Config,
) domain.SchedulePostRepository {
	return SchedulePostRepository{
		PostMessageRepository{
			db,
			config,
		},
		db,
	}
}

// SendList 配信可能なリストを取得
func (repo SchedulePostRepository) SendList(at domain.MessageSendedAt) ([]domain.SchedulePost, error) {
	models := []PostMessage{}

	if err := repo.db.GormDB.Where("message_type = ? AND sended = ? AND reservation_at <= ? AND post_messages.active = ?", domain.MessageTypeSchedulePost, false, at.Value(), true).Joins("Bot").Joins("ScheduleTiming").Find(&models).Error; err != nil {
		return []domain.SchedulePost{}, err
	}

	list := make([]domain.SchedulePost, len(models))

	for i, model := range models {
		list[i] = CreateSchedulePostEntityFromModel(model)
	}

	return list, nil
}

// Store メッセージの追加
func (repo SchedulePostRepository) Store(entity domain.SchedulePost) (vo domain.PostMessageID, err error) {
	model := PostMessage{}

	model.Message = entity.Message().Value()
	model.MessageType = entity.MessageType()
	model.BotID = entity.Bot().ID().Value()
	model.Sended = sql.NullBool{Bool: entity.IsSended(), Valid: true}
	model.Active = sql.NullBool{Bool: true, Valid: true}

	err = repo.db.GormDB.Transaction(func(tx *gorm.DB) (err error) {
		// 一旦モデルを保存
		if err = repo.db.GormDB.Create(&model).Error; err != nil {
			return err
		}

		// idを生成する
		vo, err = domain.NewPostMessageID(model.ID)
		if err != nil {
			return err
		}

		timingModel := ScheduleTiming{}
		timingModel.PostMessageID = model.ID
		timingModel.ReservationAt = entity.ReservationAt().Value()

		// 予約日時を保存
		if err = repo.db.GormDB.Create(&timingModel).Error; err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return vo, err
	}

	return vo, err
}

// Update メッセージの更新＆配信済みメッセージがあればそれも保存（予約配信なので）
func (repo SchedulePostRepository) Update(entity domain.SchedulePost) error {
	model := PostMessage{}

	model.ID = entity.ID().Value()
	model.Message = entity.Message().Value()
	model.ScheduleTiming = ScheduleTiming{
		PostMessageID: entity.ID().Value(),
		ReservationAt: entity.ReservationAt().Value(),
	}
	model.Sended = sql.NullBool{Bool: entity.IsSended(), Valid: true}

	err := repo.db.GormDB.Transaction(func(tx *gorm.DB) (err error) {
		err = repo.db.GormDB.Where("post_message_id = ?", entity.ID().Value()).Delete(&ScheduleTiming{}).Error

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

// FindByID PostMessageIDからSchedulePostを取得する
func (repo SchedulePostRepository) FindByID(id domain.PostMessageID) (entity domain.SchedulePost, err error) {
	model := PostMessage{}

	if err = repo.db.GormDB.Where("id = ? AND message_type = ?", id.Value(), domain.MessageTypeSchedulePost).First(&model).Error; err != nil {
		return entity, err
	}

	return CreateSchedulePostEntityFromModel(model), err
}

// CreateSchedulePostEntityFromModel PostMessageからEntityを生成する
func CreateSchedulePostEntityFromModel(model PostMessage) domain.SchedulePost {
	return domain.NewSchedulePost(
		model.ID,
		model.Message,
		model.ScheduleTiming.ReservationAt,
		CreateBotEntityFromModel(model.Bot),
		[]domain.SentMessage{},
		model.Sended.Bool,
	)
}

// --- DiscordMessageAdapter ---

// DiscordMessageAdapter
type DiscordMessageAdapter struct {
	upload UploadAdapter
	config core.Config
	path   core.Path
}

// NewDiscordMessageAdapter コンストラクタ
func NewDiscordMessageAdapter(
	upload UploadAdapter,
	config core.Config,
	path core.Path,
) domain.DiscordMessageAdapter {
	return DiscordMessageAdapter{
		upload,
		config,
		path,
	}
}

// SendMessage メッセージ配信リクエスト処理
func (ap DiscordMessageAdapter) SendMessage(bot domain.Bot, message domain.Message) error {

	var avatar string
	if bot.Atatar().IsEmpty() {
		// コンフィグからブランクアバターを取得する
		empty := ap.config.LoadToValueString(
			"setting",
			"resource.empty_avatar",
			"",
		)
		if empty != "" {
			avatar = ap.path.StaticURL(empty)
		}
	} else {
		avatar = ap.upload.ToURL(bot.Atatar().Value())
	}

	req := DiscordMessage{
		UserName:  bot.Name().Value(),
		AtatarURL: avatar,
		Content:   message.Value(),
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", bot.Webhook().Value(), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == 204 {
		return nil
	}

	return errors.New(fmt.Sprintf("%#v", response))
}

// DiscordMessage json変換のための構造体
type DiscordMessage struct {
	UserName  string `json:"username"`
	AtatarURL string `json:"avatar_url"`
	Content   string `json:"content"`
}

// --- PostMessage ---

// PostMessage Gormモデル
type PostMessage struct {
	gorm.Model
	Message        string             `gorm:"type:text;not null"`
	MessageType    domain.MessageType `gorm:"type:varchar(30);index"`
	BotID          uint               `gorm:"index;not null;constraint:OnDelete:CASCADE;"`
	Bot            Bot                `gorm:"constraint:OnDelete:CASCADE;"`
	SentMessages   []SentMessage      `gorm:"constraint:OnDelete:CASCADE;"`
	ScheduleTiming ScheduleTiming     `gorm:"constraint:OnDelete:CASCADE;"`
	RegularTimings []RegularTiming    `gorm:"constraint:OnDelete:CASCADE;"`
	Sended         sql.NullBool       `gorm:"type:boolean;index"`
	Active         sql.NullBool       `gorm:"type:boolean;index"`
}

// --- SentMessage ---

// SentMessage Gormモデル
type SentMessage struct {
	ID            uint   `gorm:"primarykey"`
	Message       string `gorm:"type:text;not null"`
	PostMessageID uint   `gorm:"index;not null"`
	PostMessage   PostMessage
	SendedAt      time.Time
}

// --- ScheduleTiming ---

// ScheduleTiming Gormモデル
type ScheduleTiming struct {
	ID            uint `gorm:"primarykey"`
	PostMessageID uint `gorm:"index;not null"`
	ReservationAt time.Time
}
