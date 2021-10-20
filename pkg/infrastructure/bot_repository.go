package infrastructure

import (
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
	model.Active = entity.IsActive()

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
	model.Atatar = entity.Atatar().Value()
	model.Webhook = entity.Webhook().Value()
	model.Active = entity.IsActive()

	return repo.db.GormDB.Save(&model).Error
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
	return repo.db.GormDB.Where("id = ?", id.Value()).Delete(&Bot{}).Error
}

// ExistsByWebhook BotDiscordWebhooklが重複しているBotがあるか
func (repo BotRepository) ExistsByWebhook(webhook domain.BotDiscordWebhook) (bool, error) {
	count := int64(0)
	err := repo.db.GormDB.Model(&Bot{}).
		Where("webhook = ?", webhook.Value()).
		Count(&count).
		Error

	return (count > 0), err
}

// ExistsByIDWebhook 指定したBotIDを除きBotDiscordWebhookが重複しているBotがあるか
func (repo BotRepository) ExistsByIDWebhook(id domain.BotID, webhook domain.BotDiscordWebhook) (bool, error) {
	count := int64(0)
	err := repo.db.GormDB.Model(&Bot{}).
		Where("id <> ? AND webhook = ?", id.Value(), webhook.Value()).
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
		model.Active,
	)
}

// --- Bot ---

// Bot Gormモデル
type Bot struct {
	gorm.Model
	Name    string `gorm:"type:varchar(191);index;not null"`
	Atatar  string `gorm:"type:varchar(191);index"`
	Webhook string `gorm:"type:text"`
	Active  bool   `gorm:"index"`
}