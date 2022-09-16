package configs

import "os"

type DiscordConfig struct {
	APIKey       string
	StateEnabled bool
}

func NewDiscordConfig() *DiscordConfig {
	return &DiscordConfig{
		APIKey:       os.Getenv("DISCORD_BOT_KEY"),
		StateEnabled: true,
	}
}
