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

	// message
	fx.Provide(NewPostMessageRepository),
	fx.Provide(NewImmediatePostRepository),
	fx.Provide(NewPostMessageQuery),
	fx.Provide(NewSentMessageQuery),
	fx.Provide(NewDiscordMessageAdapter),

	// adapter
	fx.Provide(NewPublicStorageUploadAdapter),
)
