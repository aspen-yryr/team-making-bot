package team_maker

import (
	"math"
)

type Player struct {
	DiscordId string
	GameId    string
}

type Scoring interface {
	GetScore(player *Player) (score float64)
}

type BaseTeamMaker struct {
	GetScore func(player *Player) (score float64)
}

func (tm *BaseTeamMaker) MakeTeam(players []*Player) [][]*Player {
	// sort players
	for i := 0; i < len(players); i++ {
		// for i, _ := range players {
		for j := i; j < len(players); j++ {
			if tm.getScore(players[i]) < tm.getScore(players[j]) {
				temp := players[i]
				players[i] = players[j]
				players[j] = temp
			}
		}
	}

	// divide 2team
	teamA := []*Player{}
	teamB := []*Player{}
	scoreA := 0.0
	scoreB := 0.0
	for _, p := range players {
		afterA := scoreA + tm.getScore(p)
		afterB := scoreB + tm.getScore(p)
		if math.Abs(afterA-scoreB) > math.Abs(scoreA-afterB) {
			teamB = append(teamB, p)
			scoreB += tm.getScore(p)
		} else {
			teamA = append(teamA, p)
			scoreA += tm.getScore(p)
		}
	}
	return [][]*Player{teamA, teamB}
}

//TODO: cache
func (tm *BaseTeamMaker) getScore(p *Player) (score float64) {
	return tm.GetScore(p)
}

func NewTeamMaker(s Scoring) *BaseTeamMaker {
	return &BaseTeamMaker{GetScore: s.GetScore}
}
