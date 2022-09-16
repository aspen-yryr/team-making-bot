package configs

import (
	"flag"
)

type BotConfig struct {
	Greet bool
}

func NewBotConfig() *BotConfig {
	return &BotConfig{
		Greet: *flag.Bool("greet", false, "If true, bot greet to channel on activated"),
	}
}
