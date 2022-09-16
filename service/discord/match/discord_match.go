package match

import (
	"context"
	"errors"
	"sync"

	"github.com/aspen-yryr/team-making-bot/internal/constants"
	"github.com/aspen-yryr/team-making-bot/pkg/dg_wrap"
	matchpb "github.com/aspen-yryr/team-making-bot/proto/match"
	dg "github.com/bwmarrin/discordgo"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var du *dg_wrap.Utility
var errs = constants.Errs

type Status int

const (
	StateVCh1Setting Status = iota
	StateVCh2Setting
	StateTeamPreview
)

// TODO: Make service
type DiscordUser struct {
	DgUser    *dg.User
	MatchUser *matchpb.User
}

type DiscordMatch struct {
	Owner              *dg.User
	tch                *dg.Channel
	Team1VCh           *dg.Channel
	Team2VCh           *dg.Channel
	MatchID            int32
	RecommendedChannel *dg.Channel
	ListeningMessage   *dg.Message
	Status             Status
}

//TODO: More strictly state management or more get method to nil check
func (m *DiscordMatch) GetRecommendedChannel() (*dg.Channel, error) {
	if m.RecommendedChannel == nil {
		return nil, errs.Unknown
	}
	return m.RecommendedChannel, nil
}

func (m *DiscordMatch) GetGuildId() string {
	return m.tch.GuildID
}

type DiscordMatchService struct {
	discordMatches []*DiscordMatch
	dcUsers        []*DiscordUser
	tchMutex       sync.RWMutex
	vchMutex       sync.RWMutex
	matchCli       matchpb.MatchSvcClient
}

func NewDiscordMatchService() *DiscordMatchService {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorf("did not connect: %v", err)
	}

	return &DiscordMatchService{
		discordMatches: []*DiscordMatch{},
		tchMutex:       sync.RWMutex{},
		vchMutex:       sync.RWMutex{},
		matchCli:       matchpb.NewMatchSvcClient(conn),
	}
}

func (s *DiscordMatchService) Create(tch *dg.Channel, owner *dg.User) (*DiscordMatch, error) {
	s.tchMutex.Lock()
	defer s.tchMutex.Unlock()

	if isContain(tch.ID, du.Channels2IDs(s.getUsingTCh())) {
		return nil, errs.MatchAlreadyStarted
	}

	user, err := s.findOrCreateUser(owner)
	if err != nil {
		return nil, err
	}

	match, err := s.matchCli.Create(
		context.TODO(),
		&matchpb.CreateMatchRequest{
			Owner: &matchpb.User{
				Id:   user.MatchUser.Id,
				Name: owner.Username,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	dMatch := &DiscordMatch{
		Owner:   owner,
		tch:     tch,
		Status:  StateVCh1Setting,
		MatchID: match.Match.Id,
	}
	s.discordMatches = append(s.discordMatches, dMatch)
	return dMatch, nil
}

func (s *DiscordMatchService) AppendMembers(tchID string, dg_users []*dg.User) error {
	match, err := s.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	members := []*matchpb.User{}
	for _, u := range dg_users {
		du, err := s.findOrCreateUser(u)
		if err != nil {
			return err
		}
		members = append(members, du.MatchUser)
	}

	_, err = s.matchCli.AppendMembers(
		context.TODO(),
		&matchpb.AppendMemberRequest{
			MatchId: match.MatchID,
			Members: members,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *DiscordMatchService) Remove(tchID string) error {
	s.tchMutex.Lock()
	defer s.tchMutex.Unlock()
	s.vchMutex.Lock()
	defer s.vchMutex.Unlock()

	for i, mt := range s.discordMatches {
		if mt.tch.ID == tchID {
			s.discordMatches[i] = s.discordMatches[len(s.discordMatches)-1]
			s.discordMatches[len(s.discordMatches)-1] = nil
			s.discordMatches = s.discordMatches[:len(s.discordMatches)-1]
			return nil
		}
	}
	return errs.MatchNotFound
}

func (s *DiscordMatchService) FilterAvailableVCh(chs []*dg.Channel) []*dg.Channel {
	vchs := du.FilterChannelsByType(chs, dg.ChannelTypeGuildVoice)
	if len(vchs) == 0 {
		return []*dg.Channel{}
	}

	availableVChs := []*dg.Channel{}
	for _, vch := range vchs {
		if !isContain(vch.ID, du.Channels2IDs(s.getUsingVCh())) {
			availableVChs = append(availableVChs, vch)
		}
	}
	return availableVChs
}

func (s *DiscordMatchService) SetVCh(tchID string, vch *dg.Channel, team string) error {
	match, err := s.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	s.vchMutex.Lock()
	defer s.vchMutex.Unlock()

	if isContain(vch.ID, du.Channels2IDs(s.getUsingVCh())) {
		return errs.ConflictVCh
	}

	if team == "Team1" {
		match.Team1VCh = vch
		match.Status = StateVCh2Setting
		match.RecommendedChannel = nil
		return nil
	} else if team == "Team2" {
		match.Team2VCh = vch
		match.Status = StateTeamPreview
		match.RecommendedChannel = nil
		return nil
	}
	return errs.InvalidTeam
}

func (s *DiscordMatchService) Shuffle(tchID string) error {
	match, err := s.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	_, err = s.matchCli.Shuffle(
		context.TODO(),
		&matchpb.ShuffleRequest{
			MatchId: match.MatchID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *DiscordMatchService) GetTeam(tchID string) (Team1UserIDs []string, Team2UserIDs []string, err error) {
	dMatch, err := s.GetMatchByTChID(tchID)
	if err != nil {
		return nil, nil, nil
	}

	match, err := s.matchCli.Find(
		context.TODO(),
		&matchpb.FindRequest{
			MatchId: dMatch.MatchID,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	listId := func(us []*matchpb.User) ([]string, error) {
		ret := []string{}
		for _, u := range us {
			du, err := s.findUserByMatchUserId(u.Id)
			if err != nil {
				return nil, err
			}
			ret = append(ret, du.DgUser.ID)
		}
		return ret, nil
	}

	tm1, err := listId(match.Match.Team1.Players)
	if err != nil {
		return nil, nil, err
	}
	tm2, err := listId(match.Match.Team2.Players)
	if err != nil {
		return nil, nil, err
	}

	return tm1, tm2, nil
}

// use cache if we need more performance (not map)
func (s *DiscordMatchService) GetMatchByTChID(tchID string) (*DiscordMatch, error) {
	for _, m := range s.discordMatches {
		if m.tch.ID == tchID {
			return m, nil
		}
	}
	return nil, errs.MatchNotFound
}

func (s *DiscordMatchService) SetListeningMessage(tchID string, msg *dg.Message) error {
	match, err := s.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}
	match.ListeningMessage = msg
	return nil
}

func (s *DiscordMatchService) IsListeningMessage(tchID, msgID string) bool {
	match, err := s.GetMatchByTChID(tchID)
	if err != nil {
		return false
	}
	if match.ListeningMessage.ID == msgID {
		return true
	}
	return false
}

// use cache if we need more performance(not map)
func (s *DiscordMatchService) getUsingTCh() []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range s.discordMatches {
		if m.tch != nil {
			using = append(using, m.tch)
		}
	}
	return using
}

func (s *DiscordMatchService) getUsingVCh() []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range s.discordMatches {
		if m.Team1VCh != nil {
			using = append(using, m.Team1VCh)
		}
		if m.Team2VCh != nil {
			using = append(using, m.Team2VCh)
		}
	}
	return using
}

func (s *DiscordMatchService) findOrCreateUser(user *dg.User) (*DiscordUser, error) {
	for _, u := range s.dcUsers {
		if u.DgUser.ID == user.ID {
			return u, nil
		}
	}

	res, err := s.matchCli.CreateUser(
		context.TODO(),
		&matchpb.CreateUserRequest{
			Name: user.Username,
		},
	)
	if err != nil {
		return nil, err
	}

	ret := &DiscordUser{
		DgUser:    user,
		MatchUser: res.User,
	}

	s.dcUsers = append(s.dcUsers, ret)
	return ret, nil
}

func (s *DiscordMatchService) findUserByMatchUserId(id int32) (*DiscordUser, error) {
	for _, u := range s.dcUsers {
		if u.MatchUser.Id == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func isContain(s string, list []string) bool {
	for _, l := range list {
		if s == l {
			return true
		}
	}
	return false
}
