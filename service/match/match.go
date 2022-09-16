package match

import (
	"context"
	"errors"
	"strconv"

	tm "github.com/aspen-yryr/team-making-bot/pkg/team-maker"
	matchpb "github.com/aspen-yryr/team-making-bot/proto/match"
	"github.com/aspen-yryr/team-making-bot/service/match/user"
)

type MatchRepository interface {
	Create() (*Match, error)
	Get(mt *Match) (*Match, error)
}

func toPbUsers(users []*user.User) []*matchpb.User {
	ret := []*matchpb.User{}
	for _, u := range users {
		ret = append(ret, &matchpb.User{
			Id:   u.ID,
			Name: u.Name,
		})
	}
	return ret
}

type Team struct {
	ID      int32
	Players []*user.User
}

func (t *Team) toPb() *matchpb.Team {
	return &matchpb.Team{
		Id:      t.ID,
		Players: toPbUsers(t.Players),
	}
}

type Match struct {
	ID      int32
	Owner   *user.User
	Team1   *Team
	Team2   *Team
	Members []*user.User
}

func (m *Match) toPb() *matchpb.Match {
	return &matchpb.Match{
		Id: m.ID,
		Owner: &matchpb.User{
			Id:   m.Owner.ID,
			Name: m.Owner.Name,
		},
		Members: toPbUsers(m.Members),
		Team1:   m.Team1.toPb(),
		Team2:   m.Team2.toPb(),
	}
}

func NewMatch(id int32, owner *user.User, tm1, tm2 *Team) *Match {
	return &Match{
		ID:      id,
		Owner:   owner,
		Team1:   tm1,
		Team2:   tm2,
		Members: []*user.User{},
	}
}

func (m *Match) AppendMember(player *user.User) {
	for _, p := range m.Members {
		if player.Is(p) {
			return
		}
	}
	m.Members = append(m.Members, player)
}

type MatchService struct {
	matches      []*Match
	users        []*user.User
	match_id_max int32
	user_id_max  int32
	team_id_max  int32
	matchpb.UnimplementedMatchSvcServer
}

func NewMatchService() matchpb.MatchSvcServer {
	// TODO: grpc client
	return &MatchService{
		matches:      []*Match{},
		users:        []*user.User{},
		match_id_max: 0,
		user_id_max:  0,
		team_id_max:  0,
	}
}

func (s *MatchService) CreateUser(_ context.Context, req *matchpb.CreateUserRequest) (*matchpb.CreateUserResponse, error) {
	u := user.New(s.user_id_max, req.Name)
	s.user_id_max++
	s.users = append(s.users, u)
	return &matchpb.CreateUserResponse{
		User: &matchpb.User{
			Id:   u.ID,
			Name: u.Name,
		},
	}, nil
}

func (s *MatchService) Create(_ context.Context, req *matchpb.CreateMatchRequest) (*matchpb.CreateMatchResponse, error) {
	for _, mt := range s.matches {
		if mt.Owner.ID == req.Owner.Id {
			return nil, errors.New("owner already has match")
		}
	}

	match := NewMatch(
		s.match_id_max,
		&user.User{
			ID:   req.Owner.Id,
			Name: req.Owner.Name,
		},
		s.CreateTeam(),
		s.CreateTeam(),
	)
	s.match_id_max++
	s.matches = append(s.matches, match)
	return &matchpb.CreateMatchResponse{
		Match: match.toPb(),
	}, nil
}

func (s *MatchService) Find(_ context.Context, req *matchpb.FindRequest) (*matchpb.FindResponse, error) {
	match, err := s.findByID(req.MatchId)
	if err != nil {
		return nil, err
	}

	return &matchpb.FindResponse{
		Match: match.toPb(),
	}, nil
}

func (s *MatchService) AppendMembers(_ context.Context, req *matchpb.AppendMemberRequest) (*matchpb.Match, error) {
	match, err := s.findByID(req.MatchId)
	if err != nil {
		return nil, err
	}

	isIn := func(id int32, users []*user.User) bool {
		for _, u := range users {
			if u.ID == id {
				return true
			}
		}
		return false
	}
	for _, mm := range req.Members {
		if isIn(mm.Id, match.Members) {
			continue
		}
		match.Members = append(match.Members, user.New(mm.Id, mm.Name))
	}

	return match.toPb(), nil
}

func (s *MatchService) Shuffle(_ context.Context, req *matchpb.ShuffleRequest) (*matchpb.ShuffleResponse, error) {
	match, err := s.findByID(req.MatchId)
	if err != nil {
		return nil, err
	}

	// TODO: Fix TeamMaker IF
	i2a := func(i []int32) []string {
		s := []string{}
		for _, ii := range i {
			s = append(s, strconv.Itoa(int(ii)))
		}
		return s
	}
	a2i := func(s []string) ([]int32, error) {
		i := []int32{}
		for _, ss := range s {
			ii, err := strconv.Atoi(ss)
			if err != nil {
				return nil, err
			}

			i = append(i, int32(ii))
		}
		return i, nil
	}

	rtm := tm.NewRandomTeamMaker(i2a(user.Ids(match.Members)))
	teams, err := rtm.MakeTeam()
	if err != nil {
		return nil, err
	}

	match.Team1.Players = []*user.User{}
	match.Team2.Players = []*user.User{}
	for _, mem := range match.Members {
		tm1, err := a2i(teams[0])
		if err != nil {
			return nil, err
		}
		for _, p := range tm1 {
			if mem.ID == p {
				match.Team1.Players = append(match.Team1.Players, mem)
			}
		}

		tm2, err := a2i(teams[1])
		if err != nil {
			return nil, err
		}
		for _, p := range tm2 {
			if mem.ID == p {
				match.Team2.Players = append(match.Team2.Players, mem)
			}
		}
	}
	return &matchpb.ShuffleResponse{}, nil
}

func (s *MatchService) CreateTeam() *Team {
	tm := &Team{
		ID:      s.team_id_max,
		Players: []*user.User{},
	}
	s.team_id_max++
	return tm
}

func (s *MatchService) findByID(id int32) (*Match, error) {
	for _, mt := range s.matches {
		if mt.ID == id {
			return mt, nil
		}
	}
	return nil, errors.New("match not found")
}
