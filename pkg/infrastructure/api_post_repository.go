package infrastructure

import (
	"database/sql"

	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/domain"
	"gorm.io/gorm"
)

// --- ApiPostRepository ---

// ApiPostRepository 予約配信Entityの永続化
type ApiPostRepository struct {
	PostMessageRepository
	db core.Database
}

// NewApiPostRepository コンストラクタ
func NewApiPostRepository(
	db core.Database,
	config core.Config,
) domain.ApiPostRepository {
	return ApiPostRepository{
		PostMessageRepository{
			db,
			config,
		},
		db,
	}
}

// Store メッセージの追加
func (repo ApiPostRepository) Store(entity domain.ApiPost) (vo domain.PostMessageID, err error) {
	model := PostMessage{}

	model.Message = entity.Message().Value()
	model.MessageType = entity.MessageType()
	model.BotID = entity.Bot().ID().Value()
	model.Sended = sql.NullBool{Bool: true, Valid: true}
	model.Active = sql.NullBool{Bool: true, Valid: true}

	if err = repo.db.GormDB.Transaction(func(tx *gorm.DB) (err error) {
		// 一旦モデルを保存
		if err = tx.Create(&model).Error; err != nil {
			return err
		}

		apiKeyModel := MessageApiKey{
			PostMessageID: model.ID,
			Key:           entity.ApiKey().Value(),
		}

		// ApiKeyを保存
		return tx.Create(&apiKeyModel).Error
	}); err != nil {
		return vo, err
	}

	// idを生成する
	return domain.NewPostMessageID(model.ID)
}

// Update メッセージの更新＆配信済みメッセージがあればそれも保存（Api配信なので）
func (repo ApiPostRepository) Update(entity domain.ApiPost) error {
	model := PostMessage{}

	model.ID = entity.ID().Value()
	model.Message = entity.Message().Value()

	err := repo.db.GormDB.Transaction(func(tx *gorm.DB) (err error) {
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

// FindByID PostMessageIDからApiPostを取得する
func (repo ApiPostRepository) FindByID(id domain.PostMessageID) (entity domain.ApiPost, err error) {
	model := PostMessage{}

	if err = repo.db.GormDB.Where("id = ? AND message_type = ?", id.Value(), domain.MessageTypeApiPost).Preload("Bot").Preload("ApiKey").First(&model).Error; err != nil {
		return entity, err
	}

	return CreateApiPostEntityFromModel(model), err
}

// FindByApiKey MessageApiKeyからApiPostを取得する
func (repo ApiPostRepository) FindByApiKey(apiKey domain.MessageApiKey) (entity domain.ApiPost, err error) {
	model := PostMessage{}

	if err = repo.db.GormDB.Where("message_type = ?", domain.MessageTypeApiPost).Preload("Bot").Preload("ApiKey").Joins(
		"LEFT JOIN message_api_keys ON message_api_keys.post_message_id = post_messages.id AND message_api_keys.key = ?",
		apiKey.Value(),
	).First(&model).Error; err != nil {
		return entity, err
	}

	return CreateApiPostEntityFromModel(model), err
}

// CreateApiPostEntityFromModel PostMessageからEntityを生成する
func CreateApiPostEntityFromModel(model PostMessage) domain.ApiPost {
	return domain.NewApiPost(
		model.ID,
		model.Message,
		model.ApiKey.Key,
		CreateBotEntityFromModel(model.Bot),
		[]domain.SentMessage{},
	)
}

// --- MessageApiKey ---

// MessageApiKey Gormモデル
type MessageApiKey struct {
	ID            uint   `gorm:"primarykey"`
	PostMessageID uint   `gorm:"index;not null;constraint:OnDelete:CASCADE;"`
	Key           string `gorm:"type:varchar(100);uniqueIndex;"`
}
