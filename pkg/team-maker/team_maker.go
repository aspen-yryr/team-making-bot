package team_maker

import "math"

type Scorer interface {
	GetScore(id string) float64
}

type TeamMaker interface {
	MakeTeam() [][]string
}

type BaseTeamMaker struct {
	DiscordIds []string
	Scorer
}

func (tm BaseTeamMaker) MakeTeam() [][]string {
	// sort players
	for i := 0; i < len(tm.DiscordIds); i++ {
		// for i, _ := range tm.DiscordIds {
		for j := i; j < len(tm.DiscordIds); j++ {
			if tm.getScore(tm.DiscordIds[i]) < tm.getScore(tm.DiscordIds[j]) {
				temp := tm.DiscordIds[i]
				tm.DiscordIds[i] = tm.DiscordIds[j]
				tm.DiscordIds[j] = temp
			}
		}
	}

	// divide 2team
	teamA := []string{}
	teamB := []string{}
	scoreA := 0.0
	scoreB := 0.0
	for _, p := range tm.DiscordIds {
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
	return [][]string{teamA, teamB}
}

//TODO: cache
func (tm BaseTeamMaker) getScore(id string) (score float64) {
	return tm.GetScore(id)
}

func NewTeamMaker(discordIDs []string, s Scorer) *BaseTeamMaker {
	tm := &BaseTeamMaker{
		DiscordIds: discordIDs,
		Scorer:     s,
	}
	return tm
}
