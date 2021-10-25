package infrastructure

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	// admin
	fx.Provide(NewAdminAuthContext),
	fx.Provide(NewAdminRepository),
	fx.Provide(NewAdminQuery),

	// bot
	fx.Provide(NewBotAtatarImageRepository),
	fx.Provide(NewBotRepository),
	fx.Provide(NewBotQuery),
	fx.Provide(NewDiscordWebhookCheckAdapter),

	// message
	fx.Provide(NewPostMessageRepository),
	fx.Provide(NewImmediatePostRepository),
	fx.Provide(NewSchedulePostRepository),
	fx.Provide(NewPostMessageQuery),
	fx.Provide(NewSentMessageQuery),
	fx.Provide(NewImmediatePostQuery),
	fx.Provide(NewSchedulePostQuery),
	fx.Provide(NewDiscordMessageAdapter),

	// adapter
	fx.Provide(NewPublicStorageUploadAdapter),
)
