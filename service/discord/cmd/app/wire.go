//go:build wireinject
// +build wireinject

package main

import (
	"github.com/aspen-yryr/team-making-bot/internal/configs"
	"github.com/aspen-yryr/team-making-bot/internal/constants"
	"github.com/aspen-yryr/team-making-bot/internal/database/connection"
	"github.com/aspen-yryr/team-making-bot/internal/database/migration"
	"github.com/aspen-yryr/team-making-bot/internal/database/repository"
	discord "github.com/aspen-yryr/team-making-bot/pkg/dg_wrap"
	"github.com/aspen-yryr/team-making-bot/service/discord/bot"
	"github.com/aspen-yryr/team-making-bot/service/discord/match"

	"github.com/google/wire"
)

func InitBot(greet bool) (*bot.Bot, error) {
	wire.Build(
		bot.Provider,
		configs.Provider,
		match.Provider,
		discord.Provider,
		constants.NewErrors,
		constants.NewMessages,
		migration.NewSampleMigrator,
		repository.Provider,
		connection.NewDevConnection,
	)
	return &bot.Bot{}, nil
}
