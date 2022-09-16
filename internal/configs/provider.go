package configs

import "github.com/google/wire"

var Provider = wire.NewSet(
	NewBotConfig,
	NewDiscordConfig,
)
