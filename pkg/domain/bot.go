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

// Equals VOの値が一致するか
func (vo BotName) Equals(eq BotName) bool {
	return vo.Value() == eq.Value()
}

// --- BotAtatar ValueObject ---

// BotAtatar Botのアバターパス
type BotAtatar string

// NewBotAtatar コンストラクタ
func NewBotAtatar(avatar string) BotAtatar {
	return BotAtatar(avatar)
}

// Value 値を返す
func (vo BotAtatar) Value() string {
	return string(vo)
}

// IsEmpty アバターが空か
func (vo BotAtatar) IsEmpty() bool {
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

// ValidURL ウェブフックURLの中のデータが正常か
func (vo BotDiscordWebhook) ValidURL(id string, token string) bool {
	return strings.Contains(vo.Value(), id+"/"+token)
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
	avatar  BotAtatar
	webhook BotDiscordWebhook
	active  bool
	userID  UserID
}

// NewBot コンストラクタ
func NewBot(
	id uint,
	name string,
	avatar string,
	webhook string,
	active bool,
	userID uint,
) Bot {
	return Bot{
		id: BotID{
			ID: Identity(id),
		},
		name:    BotName(name),
		avatar:  BotAtatar(avatar),
		webhook: BotDiscordWebhook(webhook),
		active:  active,
		userID: UserID{
			ID: Identity(userID),
		},
	}
}

// CreateBot Botを追加
func CreateBot(
	id uint,
	name string,
	avatar string,
	webhook string,
	active bool,
	userID UserID,
) (entity Bot, err error) {

	idVO, err := NewBotID(id)
	if err != nil {
		return entity, err
	}

	nameVO, err := NewBotName(name)
	if err != nil {
		return entity, err
	}

	avatarVO := NewBotAtatar(avatar)

	webhookVO, err := NewBotDiscordWebhook(webhook)
	if err != nil {
		return entity, err
	}

	return Bot{
		id:      idVO,
		name:    nameVO,
		avatar:  avatarVO,
		webhook: webhookVO,
		active:  active,
		userID:  userID,
	}, err
}

// Update Botを更新
func (entity *Bot) Update(
	name string,
	avatar string,
	webhook string,
	active bool,
) error {
	nameVO, err := NewBotName(name)
	if err != nil {
		return err
	}
	entity.name = nameVO

	entity.avatar = NewBotAtatar(avatar)

	webhookVO, err := NewBotDiscordWebhook(webhook)
	if err != nil {
		return err
	}
	entity.webhook = webhookVO

	entity.ChangeActive(active)

	return err
}

// Update UserのBotを更新
func (entity *Bot) UpdateUsers(
	name string,
	avatar string,
	webhook string,
	active bool,
	userID UserID,
) error {
	if !entity.UserID().Equals(userID) {
		return errors.New("対象UserのBotでは無いので更新できません")
	}

	return entity.Update(
		name,
		avatar,
		webhook,
		active,
	)
}

func (entity Bot) ID() BotID {
	return entity.id
}

func (entity Bot) Name() BotName {
	return entity.name
}

func (entity Bot) Atatar() BotAtatar {
	return entity.avatar
}

func (entity Bot) Webhook() BotDiscordWebhook {
	return entity.webhook
}

func (entity Bot) IsActive() bool {
	return entity.active
}

func (entity Bot) UserID() UserID {
	return entity.userID
}

// IsOwner User自身のBotかどうか
func (entity Bot) IsOwner(userID UserID) bool {
	return entity.UserID().Equals(userID)
}

// ChangeActive アクティブ状態を変更
func (entity *Bot) ChangeActive(active bool) {
	entity.active = active
}

// Equals Entityが同一か
func (entity Bot) Equals(eq Bot) bool {
	return entity.ID().Equals(eq.ID()) && entity.Name().Equals(eq.Name()) && entity.Webhook().Equals(eq.Webhook())
}

// --- BotAtatarImageFile ValueObject ---

const (
	BotAtatarContentTypeJpeg string = "image/jpeg"
	BotAtatarContentTypePng  string = "image/png"
	BotAtatarContentTypeGif  string = "image/png"
)

// GetBotAtatarContentTypes ボットのアバターコンテントタイプを取得
func GetBotAtatarContentTypes() []string {
	return []string{
		BotAtatarContentTypeJpeg,
		BotAtatarContentTypePng,
		BotAtatarContentTypeGif,
	}
}

// BotAtatarImageFile アバター画像ファイルのアップロードファイル
type BotAtatarImageFile struct {
	file *multipart.FileHeader
}

// NewBotAtatarImageFile コンストラクタ
func NewBotAtatarImageFile(
	file *multipart.FileHeader,
) (vo BotAtatarImageFile, err error) {
	if file == nil {
		return vo, errors.New("ファイルがありません")
	}

	checkTypes := GetBotAtatarContentTypes()

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		return vo, errors.New("Content-Typeヘッダーが見つかりません")
	}

	if funk.Contains(checkTypes, contentType) {
		return BotAtatarImageFile{
			file,
		}, err
	}

	return vo, errors.New("Content-TypeのMimeタイプが見つかりません")
}

// DotExt ドットを含めた拡張子を返す
func (vo BotAtatarImageFile) DotExt() string {
	name := vo.file.Filename
	pos := strings.LastIndex(name, ".")

	return name[pos:]
}

// Value 値を返す
func (vo BotAtatarImageFile) Value() *multipart.FileHeader {
	return vo.file
}

// --- BotAtatarImage Entity ---

// BotAtatarImage アバター画像ファイルEntity
type BotAtatarImage struct {
	id   UUID
	file BotAtatarImageFile
	path FilePath
}

// UploadBotAtatarImage アバター画像をアップロード
func UploadBotAtatarImage(
	file *multipart.FileHeader,
	directory string,
) (entity BotAtatarImage, err error) {
	fileVO, err := NewBotAtatarImageFile(file)

	return BotAtatarImage{
		id:   GenerateUUID(),
		file: fileVO,
		path: GenerateFilePath(
			directory,
		),
	}, err
}

func (entity BotAtatarImage) ID() UUID {
	return entity.id
}

func (entity BotAtatarImage) File() BotAtatarImageFile {
	return entity.file
}

func (entity BotAtatarImage) Path() FilePath {
	return entity.path
}

// Equals Entityが同一か
func (entity BotAtatarImage) Equals(eq BotAtatarImage) bool {
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
func (service BotService) IsDuplicate(bot Bot) (bool, error) {
	return service.repository.ExistsByNameAndWebhook(bot.Name(), bot.Webhook())
}

// IsDuplicate 指定のBotを除き重複しているか
func (service BotService) IsDuplicateWithoutSelf(bot Bot) (bool, error) {
	return service.repository.ExistsByIDNameAndWebhook(bot.ID(), bot.Name(), bot.Webhook())
}

// --- UserBotPolicy ---

// UserBotPolicy
type UserBotPolicy struct {
	context UserAuthContext
}

// NewUserBotPolicy コンストラクタ
func NewUserBotPolicy(
	context UserAuthContext,
) UserBotPolicy {
	return UserBotPolicy{
		context,
	}
}

// Detail Userが対象Botを閲覧できるか
func (policy UserBotPolicy) Detail(bot Bot) (ok bool, err error) {
	auth, err := policy.context.UserAuth()
	if err != nil {
		return ok, err
	}

	return bot.IsOwner(auth.ID()), err
}

// Update Userが対象BotをUpdateできるか
func (policy UserBotPolicy) Update(bot Bot) (ok bool, err error) {
	return policy.Detail(bot)
}

// Update Userが対象BotをDeleteできるか
func (policy UserBotPolicy) Delete(bot Bot) (ok bool, err error) {
	return policy.Detail(bot)
}

// --- BotAtatarImageRepository ---

// BotAtatarImageRepository アバター画像ファイルEntityの永続化
type BotAtatarImageRepository interface {
	Store(entity BotAtatarImage) (BotAtatar, error)
	Update(entity BotAtatarImage, avatar BotAtatar) (BotAtatar, error)
	Delete(avatar BotAtatar) error
}

// --- BotRepository ---

// BotRepository ボットEntityの永続化
type BotRepository interface {
	Store(entity Bot) (BotID, error)
	Update(entity Bot) error
	FindByID(id BotID) (Bot, error)
	Delete(id BotID) error
	ExistsByNameAndWebhook(name BotName, webhook BotDiscordWebhook) (bool, error)
	ExistsByIDNameAndWebhook(id BotID, name BotName, webhook BotDiscordWebhook) (bool, error)
	NextIdentity() (BotID, error)
}

// --- DiscordBotAdapter ---

// DiscordWebhookCheckAdapter Discordウェブフックをチェックするアダプター
type DiscordWebhookCheckAdapter interface {
	Check(webohook BotDiscordWebhook) (bool, error)
}
