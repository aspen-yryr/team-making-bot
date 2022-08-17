package team_maker

import "math"

type player struct {
	id    string
	score float64
}

type Scorer interface {
	GetScore(id string) (float64, error)
}

type TeamMaker interface {
	MakeTeam() [][]string
}

type BaseTeamMaker struct {
	DiscordIds []string
	Scorer
}

func (tm BaseTeamMaker) MakeTeam() ([][]string, error) {
	players := []*player{}
	for _, id := range tm.DiscordIds {
		score, err := tm.GetScore(id)
		if err != nil {
			return nil, err
		}

		players = append(players, &player{
			id:    id,
			score: score,
		})
	}

	// sort players
	for i := 0; i < len(players); i++ {
		for j := i; j < len(players); j++ {
			if players[i].score < players[j].score {
				temp := players[i]
				players[i] = players[j]
				players[j] = temp
			}
		}
	}

	// divide 2team
	teamA := []string{}
	teamB := []string{}
	scoreA := 0.0
	scoreB := 0.0
	for _, p := range players {
		afterA := scoreA + p.score
		afterB := scoreB + p.score
		if math.Abs(afterA-scoreB) > math.Abs(scoreA-afterB) {
			teamB = append(teamB, p.id)
			scoreB += p.score
		} else {
			teamA = append(teamA, p.id)
			scoreA += p.score
		}
	}
	return [][]string{teamA, teamB}, nil
}

func NewTeamMaker(discordIDs []string, s Scorer) *BaseTeamMaker {
	tm := &BaseTeamMaker{
		DiscordIds: discordIDs,
		Scorer:     s,
	}
	return tm
}
