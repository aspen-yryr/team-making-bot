package team_maker

import (
	"math/rand"
)

type randomTeamMaker struct {
	*BaseTeamMaker
}

type randomScorer struct {
}

func (tm randomScorer) GetScore(id string) (float64, error) {
	return rand.Float64(), nil
}

func NewRandomTeamMaker(discordIDs []string) *randomTeamMaker {
	return &randomTeamMaker{
		NewTeamMaker(discordIDs, randomScorer{}),
	}
}
