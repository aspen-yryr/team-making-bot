package dg_wrap

import "github.com/google/wire"

var Provider = wire.NewSet(
	NewDiscordSvc,
	NewSession,
)
