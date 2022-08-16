package team_maker

import (
	"math/rand"
)

type randomTeamMaker struct {
	*BaseTeamMaker
}

type randomScorer struct {
}

func (tm randomScorer) GetScore(id string) (score float64) {
	return rand.Float64()
}
func (tm randomTeamMaker) SampleInput(tmp string) {
	println(tmp)
}

func NewRandomTeamMaker(discordIDs []string) *randomTeamMaker {
	return &randomTeamMaker{
		NewTeamMaker(discordIDs, randomScorer{}),
	}
}
