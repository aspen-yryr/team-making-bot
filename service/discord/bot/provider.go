package bot

import (
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	NewBot,
)
