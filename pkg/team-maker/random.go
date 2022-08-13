package team_maker

import (
	"math/rand"
)

type randomTeamMaker struct {
	BaseTeamMaker
	// seed int
}

func (tm *randomTeamMaker) GetScore(player *Player) (score float64) {
	return rand.Float64()
}

func NewRandomTeamMaker() *randomTeamMaker {
	tm := &randomTeamMaker{}
	tm.BaseTeamMaker = *NewTeamMaker(tm)
	return tm
}

// func NewRandomTeamMaker() (*TeamMaker) {

// 	return &NewTeamMaker(tm.), nil
// }
