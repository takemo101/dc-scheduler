package infrastructure

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	// Admin
	fx.Provide(NewAdminAuthContext),
	fx.Provide(NewAdminRepository),
	fx.Provide(NewAdminQuery),

	// Bot
	fx.Provide(NewBotAtatarImageRepository),
	fx.Provide(NewBotRepository),
	fx.Provide(NewBotQuery),
	fx.Provide(NewDiscordWebhookCheckAdapter),

	// Message
	fx.Provide(NewPostMessageRepository),
	fx.Provide(NewImmediatePostRepository),
	fx.Provide(NewSchedulePostRepository),
	fx.Provide(NewRegularPostRepository),
	fx.Provide(NewPostMessageQuery),
	fx.Provide(NewSentMessageQuery),
	fx.Provide(NewImmediatePostQuery),
	fx.Provide(NewSchedulePostQuery),
	fx.Provide(NewRegularPostQuery),
	fx.Provide(NewDiscordMessageAdapter),

	// Adapter
	fx.Provide(NewPublicStorageUploadAdapter),
)
