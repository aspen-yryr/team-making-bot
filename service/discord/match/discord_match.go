package match

import (
	"context"
	"errors"
	"sync"

	"github.com/aspen-yryr/team-making-bot/internal/constants"
	"github.com/aspen-yryr/team-making-bot/pkg/discord"
	matchpb "github.com/aspen-yryr/team-making-bot/proto/match"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	dg "github.com/bwmarrin/discordgo"
)

var du *discord.Utility
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
	recommendedChannel *dg.Channel
	listeningMessage   *dg.Message
	status             Status
}

//TODO: More strictly state management or more get method to nil check
func (m *DiscordMatch) GetRecommendedChannel() (*dg.Channel, error) {
	if m.recommendedChannel == nil {
		return nil, errs.Unknown
	}
	return m.recommendedChannel, nil
}

func (m *DiscordMatch) GetGuildId() string {
	return m.tch.GuildID
}

type DiscordMatchService struct {
	dcMatches []*DiscordMatch
	dcUsers   []*DiscordUser
	tchMutex  sync.RWMutex
	vchMutex  sync.RWMutex
	svc       matchpb.MatchSvcClient
}

func NewDiscordMatchService() *DiscordMatchService {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glog.Errorf("did not connect: %v", err)
	}
	// defer conn.Close()

	return &DiscordMatchService{
		dcMatches: []*DiscordMatch{},
		tchMutex:  sync.RWMutex{},
		vchMutex:  sync.RWMutex{},
		svc:       matchpb.NewMatchSvcClient(conn),
	}
}

func (r *DiscordMatchService) Create(tch *dg.Channel, owner *dg.User) (*DiscordMatch, error) {
	r.tchMutex.Lock()
	defer r.tchMutex.Unlock()

	if isContain(tch.ID, du.Channels2IDs(r.getUsingTCh())) {
		return nil, errs.MatchAlreadyStarted
	}

	user, err := r.findOrCreateUser(owner)
	if err != nil {
		return nil, err
	}

	mt, err := r.svc.Create(
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

	dmt := &DiscordMatch{
		Owner:   owner,
		tch:     tch,
		status:  StateVCh1Setting,
		MatchID: mt.Match.Id,
	}
	r.dcMatches = append(r.dcMatches, dmt)
	return dmt, nil
}

func (r *DiscordMatchService) AppendMembers(tchID string, dg_users []*dg.User) error {
	mt, err := r.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	members := []*matchpb.User{}
	for _, u := range dg_users {
		du, err := r.findOrCreateUser(u)
		if err != nil {
			return err
		}
		members = append(members, du.MatchUser)
	}

	_, err = r.svc.AppendMembers(
		context.TODO(),
		&matchpb.AppendMemberRequest{
			MatchId: mt.MatchID,
			Members: members,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *DiscordMatchService) Remove(tchID string) error {
	r.tchMutex.Lock()
	defer r.tchMutex.Unlock()
	r.vchMutex.Lock()
	defer r.vchMutex.Unlock()

	for i, mt := range r.dcMatches {
		if mt.tch.ID == tchID {
			r.dcMatches[i] = r.dcMatches[len(r.dcMatches)-1]
			r.dcMatches[len(r.dcMatches)-1] = nil
			r.dcMatches = r.dcMatches[:len(r.dcMatches)-1]
			return nil
		}
	}
	return errs.MatchNotFound
}

func (r *DiscordMatchService) FilterAvailableVCh(chs []*dg.Channel) []*dg.Channel {
	vchs := du.FilterChannelsByType(chs, dg.ChannelTypeGuildVoice)
	if len(vchs) == 0 {
		return []*dg.Channel{}
	}

	availableVChs := []*dg.Channel{}
	for _, vch := range vchs {
		if !isContain(vch.ID, du.Channels2IDs(r.getUsingVCh())) {
			availableVChs = append(availableVChs, vch)
		}
	}
	return availableVChs
}

func (r *DiscordMatchService) SetVCh(tchID string, vch *dg.Channel, team string) error {
	mt, err := r.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	r.vchMutex.Lock()
	defer r.vchMutex.Unlock()

	if isContain(vch.ID, du.Channels2IDs(r.getUsingVCh())) {
		return errs.ConflictVCh
	}

	if team == "Team1" {
		mt.Team1VCh = vch
		mt.status = StateVCh2Setting
		mt.recommendedChannel = nil
		return nil
	} else if team == "Team2" {
		mt.Team2VCh = vch
		mt.status = StateTeamPreview
		mt.recommendedChannel = nil
		return nil
	}
	return errs.InvalidTeam
}

func (r *DiscordMatchService) Shuffle(tchID string) error {
	mt, err := r.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	_, err = r.svc.Shuffle(
		context.TODO(),
		&matchpb.ShuffleRequest{
			MatchId: mt.MatchID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DiscordMatchService) GetTeam(tchID string) (Team1UserIDs []string, Team2UserIDs []string, err error) {
	dmt, err := r.GetMatchByTChID(tchID)
	if err != nil {
		return nil, nil, nil
	}

	mt, err := r.svc.Find(
		context.TODO(),
		&matchpb.FindRequest{
			MatchId: dmt.MatchID,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	f := func(us []*matchpb.User) ([]string, error) {
		ret := []string{}
		for _, u := range us {
			du, err := r.findUserByMatchUserId(u.Id)
			if err != nil {
				return nil, err
			}
			ret = append(ret, du.DgUser.ID)
		}
		return ret, nil
	}

	tm1, err := f(mt.Match.Team1.Players)
	if err != nil {
		return nil, nil, err
	}
	tm2, err := f(mt.Match.Team2.Players)
	if err != nil {
		return nil, nil, err
	}

	return tm1, tm2, nil
}

// use cache if we need more performance (not map)
func (r *DiscordMatchService) GetMatchByTChID(tchID string) (*DiscordMatch, error) {
	for _, m := range r.dcMatches {
		if m.tch.ID == tchID {
			return m, nil
		}
	}
	return nil, errs.MatchNotFound
}

func (r *DiscordMatchService) GetMatchStatus(tchID string) (*Status, error) {
	mt, err := r.GetMatchByTChID(tchID)
	if err != nil {
		return nil, err
	}
	return &mt.status, nil
}

func (r *DiscordMatchService) SetRecommendedChannel(tchID string, vch *dg.Channel) error {
	mt, err := r.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}
	mt.recommendedChannel = vch
	return nil
}

func (r *DiscordMatchService) SetListeningMessage(tchID string, msg *dg.Message) error {
	mt, err := r.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}
	mt.listeningMessage = msg
	return nil
}

func (r *DiscordMatchService) IsListeningMessage(tchID, msgID string) bool {
	mt, err := r.GetMatchByTChID(tchID)
	if err != nil {
		return false
	}
	if mt.listeningMessage.ID == msgID {
		return true
	}
	return false
}

// use cache if we need more performance(not map)
func (r *DiscordMatchService) getUsingTCh() []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range r.dcMatches {
		if m.tch != nil {
			using = append(using, m.tch)
		}
	}
	return using
}

func (r *DiscordMatchService) getUsingVCh() []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range r.dcMatches {
		if m.Team1VCh != nil {
			using = append(using, m.Team1VCh)
		}
		if m.Team2VCh != nil {
			using = append(using, m.Team2VCh)
		}
	}
	return using
}

func (r *DiscordMatchService) findOrCreateUser(user *dg.User) (*DiscordUser, error) {
	for _, u := range r.dcUsers {
		if u.DgUser.ID == user.ID {
			return u, nil
		}
	}

	res, err := r.svc.CreateUser(
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

	r.dcUsers = append(r.dcUsers, ret)
	return ret, nil
}

func (r *DiscordMatchService) findUserByMatchUserId(id int32) (*DiscordUser, error) {
	for _, u := range r.dcUsers {
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
