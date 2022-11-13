package match

import (
	"team-making-bot/internal/match/player"
	tm "team-making-bot/pkg/team-maker"
)

type match struct {
	Owner *player.Player
	Team1 []*player.Player
	Team2 []*player.Player
}

func NewMatch(owner *player.Player) *match {
	return &match{
		Owner: owner,
		Team1: []*player.Player{},
		Team2: []*player.Player{},
	}
}

func (m *match) Ids() []string {
	var ids []string
	for _, p := range append(m.Team1, m.Team2...) {
		ids = append(ids, p.ID)
	}
	return ids
}

func (m *match) AppendMember(id string) {
	panic("Append Member not implemented")
}

type MatchService struct {
	list []*match
}

func NewMatchService() *MatchService {
	return &MatchService{
		list: []*match{},
	}
}

func (m *MatchService) Find(owner *player.Player) (*match, error) {
	for _, mt := range m.list {
		if mt.Owner.Is(owner) {
			return mt, nil
		}
	}
	return nil, errs.MatchNotFound
}

func (m *MatchService) Create(owner *player.Player) (*match, error) {
	_, err := m.Find(owner)
	if err != errs.MatchNotFound && err != nil {
		return nil, err
	}

	mt := NewMatch(owner)
	m.list = append(m.list, mt)
	return mt, nil
}

func (m *MatchService) Remove(owner *player.Player) error {
	for i, mt := range m.list {
		if mt.Owner.Is(owner) {
			m.list[i] = m.list[len(m.list)-1]
			m.list[len(m.list)-1] = nil
			m.list = m.list[:len(m.list)-1]
			return nil
		}
	}
	return errs.MatchNotFound
}

func (m *MatchService) Shuffle(owner *player.Player) error {
	mt, err := m.Find(owner)
	if err != nil {
		return err
	}

	Ids := mt.Ids()
	players := append(mt.Team1, mt.Team2...)

	rtm := tm.NewRandomTeamMaker(Ids)
	teams, err := rtm.MakeTeam()
	if err != nil {
		return err
	}

	mt.Team1 = []*player.Player{}
	mt.Team2 = []*player.Player{}
	for _, p := range players {
		for _, t1 := range teams[0] {
			if p.ID == t1 {
				mt.Team1 = append(mt.Team1, p)
			}
		}
		for _, t2 := range teams[1] {
			if p.ID == t2 {
				mt.Team2 = append(mt.Team2, p)
			}
		}
	}
	return nil
}
