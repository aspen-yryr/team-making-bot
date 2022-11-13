package match

import (
	"team-making-bot/internal/match/user"
	tm "team-making-bot/pkg/team-maker"
)

type Match struct {
	Owner   *user.User
	Team1   []*user.User
	Team2   []*user.User
	Players []*user.User
}

func NewMatch(owner *user.User) *Match {
	return &Match{
		Owner:   owner,
		Team1:   []*user.User{},
		Team2:   []*user.User{},
		Players: []*user.User{},
	}
}

func (m *Match) AppendMember(player *user.User) {
	for _, p := range m.Players {
		if player.Is(p) {
			return
		}
	}
	m.Players = append(m.Players, player)
}

type MatchService struct {
	list []*Match
}

func NewMatchService() *MatchService {
	return &MatchService{
		list: []*Match{},
	}
}

func (m *MatchService) Find(owner *user.User) (*Match, error) {
	for _, mt := range m.list {
		if mt.Owner.Is(owner) {
			return mt, nil
		}
	}
	return nil, errs.MatchNotFound
}

func (m *MatchService) Create(owner *user.User) (*Match, error) {
	_, err := m.Find(owner)
	if err != errs.MatchNotFound && err != nil {
		return nil, err
	}

	mt := NewMatch(owner)
	m.list = append(m.list, mt)
	return mt, nil
}

func (m *MatchService) Remove(owner *user.User) error {
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

func (m *MatchService) Shuffle(owner *user.User) error {
	mt, err := m.Find(owner)
	if err != nil {
		return err
	}

	rtm := tm.NewRandomTeamMaker(user.Ids(mt.Players))
	teams, err := rtm.MakeTeam()
	if err != nil {
		return err
	}

	mt.Team1 = []*user.User{}
	mt.Team2 = []*user.User{}
	for _, p := range mt.Players {
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
