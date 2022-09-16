//go:build wireinject
// +build wireinject

package match

import (
	"team-making-bot/pkg/database/connection"
	"team-making-bot/pkg/database/repository"
	team_maker "team-making-bot/pkg/team-maker"

	"github.com/google/wire"
)

func InitBattleRecordTeamMaker() (*team_maker.BattleRecordTeamMaker, error) {
	wire.Build(
		team_maker.NewBattleRecordTeamMaker,
		repository.NewBattleRecordRepository,
		connection.NewDevConnection,
	)
	return &team_maker.BattleRecordTeamMaker{}, nil
}
