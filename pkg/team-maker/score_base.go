package team_maker

import (
	"math/rand"
	"time"
)

type scoreBaseTeamMaker struct {
	BaseTeamMaker
	scores map[string]int
	// seed int
}

func (tm *scoreBaseTeamMaker) MakeTeams(players []*Player) ([]*Player, []*Player) {
	l := int(len(players) / 2)
	rand.Seed(int64(time.Now().UnixNano()))
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})
	return players[:l], players[l:]
}

func NewScoreBaseTeamMaker() (*scoreBaseTeamMaker, error) {
	return &scoreBaseTeamMaker{}, nil
}
