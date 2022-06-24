package teammaker

import (
	"math/rand"
	"time"
)

type Player struct {
	DiscordId string
	GameId    string
	Name      string
}

type TeamMaker interface {
	MakeTeam(players []Player) (teamA []Player, teamB []Player)
}

type randomTeamMaker struct {
	// seed int
}

func (tm *randomTeamMaker) MakeTeam(players []*Player) ([]*Player, []*Player) {
	l := int(len(players) / 2)
	rand.Seed(int64(time.Now().UnixNano()))
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})
	return players[:l], players[l:]
}

func NewRandomTeamMaker() (*randomTeamMaker, error) {
	return &randomTeamMaker{}, nil
}
