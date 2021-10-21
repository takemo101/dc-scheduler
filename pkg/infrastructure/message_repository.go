package infrastructure

import (
	"bytes"
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
	return repo.db.GormDB.Where("id = ?", id.Value()).Delete(&PostMessage{}).Error
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

// Store メッセージの追加＆送信済みメッセージがあればそれも保存（即時配信なので）
func (repo ImmediatePostRepository) Store(entity domain.ImmediatePost) (vo domain.PostMessageID, err error) {
	model := PostMessage{}

	model.Message = entity.Message().Value()
	model.MessageType = entity.MessageType()
	model.BotID = entity.Bot().ID().Value()

	// 一旦モデルを保存
	if err = repo.db.GormDB.Create(&model).Error; err != nil {
		return vo, err
	}

	// idを生成する
	vo, err = domain.NewPostMessageID(model.ID)
	if err != nil {
		return vo, err
	}

	// 送信済みであればメッセージを保存
	if entity.IsSended() {
		for _, sent := range entity.SentMessages() {
			err = repo.SaveSendedMessage(vo, sent)
			if err != nil {
				return vo, err
			}
		}
	}

	return vo, err
}

// FindByID PostMessageIDからImmediatePostを取得する
func (repo ImmediatePostRepository) FindByID(id domain.PostMessageID) (entity domain.ImmediatePost, err error) {
	model := PostMessage{}

	if err = repo.db.GormDB.Where("id = ?", id.Value()).Preload("Bot").First(&model).Error; err != nil {
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

// SendMessage メッセージ送信リクエスト処理
func (ap DiscordMessageAdapter) SendMessage(bot domain.Bot, message domain.Message) error {

	var avatar string
	if bot.Atatar().IsEmpty() {
		avatar = ap.upload.ToURL(bot.Atatar().Value())
	} else {
		// コンフィグからブランクアバターを取得する
		empty := ap.config.LoadToValueString(
			"setting",
			"resource.empty_avatar",
			"",
		)
		if empty != "" {
			avatar = ap.path.StaticURL(empty)
		}
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
	BotID          uint               `gorm:"index;not null"`
	Bot            Bot                `gorm:"constraint:OnDelete:CASCADE;"`
	SentMessages   []SentMessage      `gorm:"constraint:OnDelete:CASCADE;"`
	ScheduleTiming ScheduleTiming     `gorm:"constraint:OnDelete:CASCADE;"`
}

// --- SentMessage ---

// SentMessage Gormモデル
type SentMessage struct {
	ID            uint   `gorm:"primarykey"`
	Message       string `gorm:"type:text;not null"`
	PostMessageID uint   `gorm:"index;not null"`
	SendedAt      time.Time
}

// --- ScheduleTiming ---

// ScheduleTiming Gormモデル
type ScheduleTiming struct {
	ID            uint `gorm:"primarykey"`
	PostMessageID uint `gorm:"index;not null"`
	ReservationAt time.Time
}
