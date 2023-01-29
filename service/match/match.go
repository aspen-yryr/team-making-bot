package match

import (
	"context"
	"errors"
	"strconv"

	tm "github.com/aspen-yryr/team-making-bot/pkg/team-maker"
	matchpb "github.com/aspen-yryr/team-making-bot/proto/match"
	"github.com/aspen-yryr/team-making-bot/service/match/user"
	"github.com/golang/glog"
)

func toPbUsers(users []*user.User) []*matchpb.User {
	upb := []*matchpb.User{}
	for _, u := range users {
		upb = append(upb, &matchpb.User{
			Id:   u.ID,
			Name: u.Name,
		})
	}
	return upb
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
	return &MatchService{
		matches:      []*Match{},
		users:        []*user.User{},
		match_id_max: 0,
		user_id_max:  0,
		team_id_max:  0,
	}
}

func (m *MatchService) CreateUser(_ context.Context, req *matchpb.CreateUserRequest) (*matchpb.CreateUserResponse, error) {
	glog.V(6).Infoln("CreateUser called")
	u := user.New(m.user_id_max, req.Name)
	m.user_id_max++
	m.users = append(m.users, u)
	return &matchpb.CreateUserResponse{
		User: &matchpb.User{
			Id:   u.ID,
			Name: u.Name,
		},
	}, nil
}

func (m *MatchService) Create(_ context.Context, req *matchpb.CreateMatchRequest) (*matchpb.CreateMatchResponse, error) {
	glog.V(6).Infoln("Create called")
	for _, mt := range m.matches {
		if mt.Owner.ID == req.Owner.Id {
			return nil, errors.New("owner already has match")
		}
	}

	mt := NewMatch(
		m.match_id_max,
		&user.User{
			ID:   req.Owner.Id,
			Name: req.Owner.Name,
		},
		m.CreateTeam(),
		m.CreateTeam(),
	)
	m.match_id_max++
	m.matches = append(m.matches, mt)
	return &matchpb.CreateMatchResponse{
		Match: mt.toPb(),
	}, nil
}

func (m *MatchService) Find(_ context.Context, req *matchpb.FindRequest) (*matchpb.FindResponse, error) {
	glog.V(6).Infoln("Find called")
	mt, err := m.findByID(req.MatchId)
	if err != nil {
		return nil, err
	}

	return &matchpb.FindResponse{
		Match: mt.toPb(),
	}, nil
}

func (m *MatchService) AppendMembers(_ context.Context, req *matchpb.AppendMemberRequest) (*matchpb.Match, error) {
	mt, err := m.findByID(req.MatchId)
	if err != nil {
		return nil, err
	}

	in := func(id int32, us []*user.User) bool {
		for _, u := range us {
			if u.ID == id {
				return true
			}
		}
		return false
	}
	for _, mm := range req.Members {
		if in(mm.Id, mt.Members) {
			continue
		}
		mt.Members = append(mt.Members, &user.User{ID: mm.Id, Name: mm.Name})
	}
	return mt.toPb(), nil
}

func (m *MatchService) Shuffle(_ context.Context, req *matchpb.ShuffleRequest) (*matchpb.ShuffleResponse, error) {
	mt, err := m.findByID(req.MatchId)
	if err != nil {
		return nil, err
	}

	// TODO: Fix TeamMaker IF
	f := func(i []int32) []string {
		s := []string{}
		for _, ii := range i {
			s = append(s, strconv.Itoa(int(ii)))
		}
		return s
	}
	f_ := func(ss []string) ([]int32, error) {
		i := []int32{}
		for _, str := range ss {
			ii, err := strconv.Atoi(str)
			if err != nil {
				return nil, err
			}
			i = append(i, int32(ii))
		}
		return i, nil
	}

	rtm := tm.NewRandomTeamMaker(f(user.Ids(mt.Members)))
	teams, err := rtm.MakeTeam()
	if err != nil {
		return nil, err
	}

	mt.Team1.Players = []*user.User{}
	mt.Team2.Players = []*user.User{}
	for _, mem := range mt.Members {
		tm1, err := f_(teams[0])
		if err != nil {
			return nil, err
		}
		for _, p := range tm1 {
			if mem.ID == p {
				mt.Team1.Players = append(mt.Team1.Players, mem)
			}
		}
		tm2, err := f_(teams[1])
		if err != nil {
			return nil, err
		}
		for _, p := range tm2 {
			if mem.ID == p {
				mt.Team2.Players = append(mt.Team2.Players, mem)
			}
		}
	}
	return &matchpb.ShuffleResponse{}, nil
}

func (m *MatchService) CreateTeam() *Team {
	tm := &Team{
		ID:      m.team_id_max,
		Players: []*user.User{},
	}
	m.team_id_max++
	return tm
}

func (m *MatchService) findByID(id int32) (*Match, error) {
	for _, mt := range m.matches {
		if mt.ID == id {
			return mt, nil
		}
	}
	return nil, errors.New("match not found")
}
