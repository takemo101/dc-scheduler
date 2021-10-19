package domain

import (
	"errors"
	"mime/multipart"
	"strings"
	"unicode/utf8"

	"github.com/thoas/go-funk"
)

// --- BotID ValueObject ---

// BotID BotのID
type BotID struct {
	ID Identity
}

// NewBotID コンストラクタ
func NewBotID(id uint) (vo BotID, err error) {

	if identity, err := NewIdentity(id); err == nil {
		return BotID{
			ID: identity,
		}, err
	}

	return vo, err
}

// Value IDの値を返す
func (vo BotID) Value() uint {
	return vo.ID.Value()
}

// Equals VOの値が一致するか
func (vo BotID) Equals(eq BotID) bool {
	return vo.Value() == eq.Value()
}

// --- BotName ValueObject ---

// BotName Botの名前
type BotName string

// NewBotName コンストラクタ
func NewBotName(name string) (vo BotName, err error) {
	length := utf8.RuneCountInString(name)

	// 32文字以上は設定できない
	if length == 0 || length > 80 {
		return vo, errors.New("BotNameは80文字以内で設定してください")
	}

	return BotName(name), err
}

// Value 値を返す
func (vo BotName) Value() string {
	return string(vo)
}

// --- BotAvator ValueObject ---

// BotAvator Botのアバターパス
type BotAvator string

// NewBotAvator コンストラクタ
func NewBotAvator(avator string) BotAvator {
	return BotAvator(avator)
}

// Value 値を返す
func (vo BotAvator) Value() string {
	return string(vo)
}

// IsEmpty アバターが空か
func (vo BotAvator) IsEmpty() bool {
	return len(vo.Value()) == 0
}

// --- BotDiscordWebhook ValueObject ---

// BotDiscordWebhook BotのウェブフックURL
type BotDiscordWebhook string

const BotDiscordWebhookURLPrefix string = "https://discord.com/api/webhooks/"

// NewBotDiscordWebhook コンストラクタ
func NewBotDiscordWebhook(webhook string) (vo BotDiscordWebhook, err error) {
	if len(webhook) == 0 || !strings.Contains(webhook, BotDiscordWebhookURLPrefix) {
		return vo, errors.New("BotDiscordWebhookのURLが不正です")
	}

	return BotDiscordWebhook(webhook), err
}

// Value 値を返す
func (vo BotDiscordWebhook) Value() string {
	return string(vo)
}

// Equals VOの値が一致するか
func (vo BotDiscordWebhook) Equals(eq BotDiscordWebhook) bool {
	return vo.Value() == eq.Value()
}

// --- Bot Entity ---

// Bot ボットEntity
type Bot struct {
	id      BotID
	name    BotName
	avator  BotAvator
	webhook BotDiscordWebhook
	active  bool
}

// NewBot コンストラクタ
func NewBot(
	id uint,
	name string,
	avator string,
	webhook string,
	active bool,
) Bot {
	return Bot{
		id: BotID{
			ID: Identity(id),
		},
		name:    BotName(name),
		avator:  BotAvator(avator),
		webhook: BotDiscordWebhook(webhook),
		active:  active,
	}
}

// CreateBot Botを追加
func CreateBot(
	id uint,
	name string,
	avator string,
	webhook string,
	active bool,
) (entity Bot, err error) {

	idVO, err := NewBotID(id)
	if err != nil {
		return entity, err
	}

	nameVO, err := NewBotName(name)
	if err != nil {
		return entity, err
	}

	avatorVO := NewBotAvator(avator)

	webhookVO, err := NewBotDiscordWebhook(webhook)
	if err != nil {
		return entity, err
	}

	return Bot{
		id:      idVO,
		name:    nameVO,
		avator:  avatorVO,
		webhook: webhookVO,
		active:  active,
	}, err
}

// Update Botを更新
func (entity *Bot) Update(
	name string,
	avator string,
	webhook string,
	active bool,
) error {
	nameVO, err := NewBotName(name)
	if err != nil {
		return err
	}
	entity.name = nameVO

	entity.avator = NewBotAvator(avator)

	webhookVO, err := NewBotDiscordWebhook(webhook)
	if err != nil {
		return err
	}
	entity.webhook = webhookVO

	entity.ChangeActive(active)

	return err
}

func (entity Bot) ID() BotID {
	return entity.id
}

func (entity Bot) Name() BotName {
	return entity.name
}

func (entity Bot) Avator() BotAvator {
	return entity.avator
}

func (entity Bot) Webhook() BotDiscordWebhook {
	return entity.webhook
}

func (entity Bot) IsActive() bool {
	return entity.active
}

// ChangeActive アクティブ状態を変更
func (entity *Bot) ChangeActive(active bool) {
	entity.active = active
}

// Equals Entityが同一か
func (entity Bot) Equals(eq Bot) bool {
	return entity.ID().Equals(eq.ID()) && entity.Webhook().Equals(eq.Webhook())
}

// --- BotAvatorImageFile ValueObject ---

const (
	BotAvatorContentTypeJpeg string = "image/jpeg"
	BotAvatorContentTypePng  string = "image/png"
	BotAvatorContentTypeGif  string = "image/png"
)

// GetBotAvatorContentTypes ボットのアバターコンテントタイプを取得
func GetBotAvatorContentTypes() []string {
	return []string{
		BotAvatorContentTypeJpeg,
		BotAvatorContentTypePng,
		BotAvatorContentTypeGif,
	}
}

// BotAvatorImageFile アバター画像ファイルのアップロードファイル
type BotAvatorImageFile struct {
	file *multipart.FileHeader
}

// NewBotAvatorImageFile コンストラクタ
func NewBotAvatorImageFile(
	file *multipart.FileHeader,
) (vo BotAvatorImageFile, err error) {
	if file == nil {
		return vo, errors.New("ファイルがありません")
	}

	checkTypes := GetBotAvatorContentTypes()

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		return vo, errors.New("Content-Typeヘッダーが見つかりません")
	}

	if funk.Contains(checkTypes, contentType) {
		return BotAvatorImageFile{
			file,
		}, err
	}

	return vo, errors.New("Content-TypeのMimeタイプが見つかりません")
}

// DotExt ドットを含めた拡張子を返す
func (vo BotAvatorImageFile) DotExt() string {
	name := vo.file.Filename
	pos := strings.LastIndex(name, ".")

	return name[pos:]
}

// Value 値を返す
func (vo BotAvatorImageFile) Value() *multipart.FileHeader {
	return vo.file
}

// --- BotAvatorImage Entity ---

// BotAvatorImage アバター画像ファイルEntity
type BotAvatorImage struct {
	id   UUID
	file BotAvatorImageFile
	path FilePath
}

// UploadBotAvatorImage アバター画像をアップロード
func UploadBotAvatorImage(
	file *multipart.FileHeader,
	directory string,
) (entity BotAvatorImage, err error) {
	fileVO, err := NewBotAvatorImageFile(file)

	return BotAvatorImage{
		id:   GenerateUUID(),
		file: fileVO,
		path: GenerateFilePath(
			directory,
		),
	}, err
}

func (entity BotAvatorImage) ID() UUID {
	return entity.id
}

func (entity BotAvatorImage) File() BotAvatorImageFile {
	return entity.file
}

func (entity BotAvatorImage) Path() FilePath {
	return entity.path
}

// Equals Entityが同一か
func (entity BotAvatorImage) Equals(eq BotAvatorImage) bool {
	return entity.ID().Equals(eq.ID())
}

// --- BotService ---

// BotService 管理者ドメインサービス
type BotService struct {
	repository BotRepository
}

// NewBotService コンストラクタ
func NewBotService(
	repository BotRepository,
) BotService {
	return BotService{
		repository,
	}
}

// IsDuplicate Botが重複しているか
func (service BotService) IsDuplicate(Bot Bot) (bool, error) {
	return service.repository.ExistsByWebhook(Bot.Webhook())
}

// IsDuplicate 指定のBotを除き重複しているか
func (service BotService) IsDuplicateWithoutSelf(Bot Bot) (bool, error) {
	return service.repository.ExistsByIDWebhook(Bot.ID(), Bot.Webhook())
}

// --- BotAvatorImageRepository ---

// BotAvatorImageRepository アバター画像ファイルEntityの永続化
type BotAvatorImageRepository interface {
	Store(entity BotAvatorImage) (BotAvator, error)
	Update(entity BotAvatorImage, avator BotAvator) (BotAvator, error)
	Delete(avator BotAvator) error
}

// --- BotRepository ---

// BotRepository ボットEntityの永続化
type BotRepository interface {
	Store(entity Bot) (BotID, error)
	Update(entity Bot) error
	FindByID(id BotID) (Bot, error)
	Delete(id BotID) error
	ExistsByWebhook(webhook BotDiscordWebhook) (bool, error)
	ExistsByIDWebhook(id BotID, webhook BotDiscordWebhook) (bool, error)
	NextIdentity() (BotID, error)
}
