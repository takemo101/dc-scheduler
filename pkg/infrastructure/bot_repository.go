package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/domain"
	"gorm.io/gorm"
)

// --- BotAtatarImageRepository ---

// BotAtatarImageRepository アバター画像ファイルEntityの永続化
type BotAtatarImageRepository struct {
	upload UploadAdapter
}

func NewBotAtatarImageRepository(
	upload UploadAdapter,
) domain.BotAtatarImageRepository {
	return BotAtatarImageRepository{
		upload,
	}
}

// Store アバター画像を保存する
func (repo BotAtatarImageRepository) Store(entity domain.BotAtatarImage) (vo domain.BotAtatar, err error) {
	// ディレクトリがなければ作成する
	if !repo.upload.Exists(entity.Path().Directory()) {
		err = repo.upload.MakeDirectory(entity.Path().Directory())
		if err != nil {
			return vo, err
		}
	}

	toPath := entity.Path().Value() + entity.File().DotExt()
	// アップロード処理
	_, err = repo.upload.Upload(
		entity.File().Value(),
		toPath,
	)
	if err != nil {
		return vo, err
	}

	return domain.NewBotAtatar(toPath), err
}

// Update アバター画像を更新する
func (repo BotAtatarImageRepository) Update(entity domain.BotAtatarImage, avatar domain.BotAtatar) (vo domain.BotAtatar, err error) {
	// 更新されたファイルをアップロードする
	vo, err = repo.Store(entity)
	if err != nil {
		return vo, err
	}

	// 前のアバターが存在すれば削除しとく
	return vo, repo.Delete(avatar)
}

// Delete アバター画像を削除する
func (repo BotAtatarImageRepository) Delete(avatar domain.BotAtatar) (err error) {
	if !avatar.IsEmpty() {
		return repo.upload.Delete(avatar.Value())
	}

	return err
}

// --- BotRepository ---

// BotRepository ボットEntityの永続化
type BotRepository struct {
	db     core.Database
	config core.Config
}

func NewBotRepository(
	db core.Database,
	config core.Config,
) domain.BotRepository {
	return BotRepository{
		db,
		config,
	}
}

// Store Botの追加
func (repo BotRepository) Store(entity domain.Bot) (vo domain.BotID, err error) {
	model := Bot{}

	model.Name = entity.Name().Value()
	model.Atatar = entity.Atatar().Value()
	model.Webhook = entity.Webhook().Value()
	model.Active = sql.NullBool{Bool: entity.IsActive(), Valid: true}
	model.UserID = entity.UserID().Value()

	if err = repo.db.GormDB.Create(&model).Error; err != nil {
		return vo, err
	}

	return domain.NewBotID(model.ID)
}

// Update Botの更新
func (repo BotRepository) Update(entity domain.Bot) error {
	model := Bot{}

	model.ID = entity.ID().Value()
	model.Name = entity.Name().Value()
	if !entity.Atatar().IsEmpty() {
		model.Atatar = entity.Atatar().Value()
	}
	model.Webhook = entity.Webhook().Value()
	model.Active = sql.NullBool{Bool: entity.IsActive(), Valid: true}

	return repo.db.GormDB.Updates(&model).Error
}

// FindByID BotIDからBotを取得する
func (repo BotRepository) FindByID(id domain.BotID) (entity domain.Bot, err error) {
	model := Bot{}

	if err = repo.db.GormDB.Where("id = ?", id.Value()).First(&model).Error; err != nil {
		return entity, err
	}

	return CreateBotEntityFromModel(model), err
}

// Delete BotIDからBotを削除する
func (repo BotRepository) Delete(id domain.BotID) error {
	return repo.db.GormDB.Where("id = ?", id.Value()).Unscoped().Delete(&Bot{}).Error
}

// ExistsByWebhook BotDiscordWebhooklが重複しているBotがあるか
func (repo BotRepository) ExistsByNameAndWebhookAndUserID(name domain.BotName, webhook domain.BotDiscordWebhook, userID domain.UserID) (bool, error) {
	count := int64(0)
	err := repo.db.GormDB.Model(&Bot{}).
		Where(
			"name = ? AND webhook = ? AND user_id = ?",
			name.Value(),
			webhook.Value(),
			userID.Value(),
		).
		Count(&count).
		Error

	return (count > 0), err
}

// ExistsByIDNameAndWebhook 指定したBotIDを除きBotDiscordWebhookが重複しているBotがあるか
func (repo BotRepository) ExistsByIDNameAndWebhookAndUserID(id domain.BotID, name domain.BotName, webhook domain.BotDiscordWebhook, userID domain.UserID) (bool, error) {
	count := int64(0)
	err := repo.db.GormDB.Model(&Bot{}).
		Where(
			"id <> ? AND name = ? AND webhook = ? AND user_id = ?",
			id.Value(),
			name.Value(),
			webhook.Value(),
			userID.Value(),
		).
		Count(&count).
		Error

	return (count > 0), err
}

// NextIdentity 次のIDを取得する
func (repo BotRepository) NextIdentity() (domain.BotID, error) {
	var max uint

	sql := GetNextIdentitySelectSQL(repo.config.DB.Type)
	repo.db.GormDB.Model(&Bot{}).Select(sql).Scan(&max)

	return domain.NewBotID(max + 1)
}

// CreateBotEntityFromModel BotからEntityを生成する
func CreateBotEntityFromModel(model Bot) domain.Bot {
	return domain.NewBot(
		model.ID,
		model.Name,
		model.Atatar,
		model.Webhook,
		model.Active.Bool,
		model.UserID,
	)
}

// ---- DiscordWebhookCheckAdapter ---

// DiscordWebhookCheckAdapter Discordウェブフックをチェックするアダプター
type DiscordWebhookCheckAdapter struct {
	//
}

// NewDiscordWebhookCheckAdapter コンストラクタ
func NewDiscordWebhookCheckAdapter(
//
) domain.DiscordWebhookCheckAdapter {
	return DiscordWebhookCheckAdapter{}
}

// Check 指定したウェブフックがアクセス可能か
func (ap DiscordWebhookCheckAdapter) Check(
	webhook domain.BotDiscordWebhook,
) (ok bool, err error) {
	response, err := http.Get(webhook.Value())
	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	bytes := []byte(body)

	var checkResponse WebhookCheckResponse
	err = json.Unmarshal(bytes, &checkResponse)
	if err != nil {
		return false, err
	}

	if response.StatusCode == 200 {
		if webhook.ValidURL(
			checkResponse.ID,
			checkResponse.Token,
		) {
			return true, err
		}

		return false, errors.New("Webookとレスポンス値が一致しません")
	}

	return ok, errors.New(fmt.Sprintf("%#v", response))
}

type WebhookCheckResponse struct {
	Type          int    `json:"type"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Avator        string `json:"avator"`
	ChannelID     string `json:"channel_id"`
	GuildID       string `json:"guild_id"`
	ApplicationID string `json:"application_id"`
	Token         string `json:"token"`
}

// --- Bot ---

// Bot Gormモデル
type Bot struct {
	gorm.Model
	Name    string       `gorm:"type:varchar(191);index;not null"`
	Atatar  string       `gorm:"type:varchar(191);index"`
	Webhook string       `gorm:"type:text"`
	Active  sql.NullBool `gorm:"type:boolean;index"`
	UserID  uint         `gorm:"index;not null"`
	User    User         `gorm:"constraint:OnDelete:CASCADE;"`
}
