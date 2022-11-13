package match

import (
	"sync"
	"team-making-bot/internal/constants"
	"team-making-bot/internal/match/user"
	"team-making-bot/pkg/discord"

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

type DiscordMatch struct {
	// TODO: dg.channel getter for nil check
	Owner              *dg.User
	tch                *dg.Channel
	Team1VCh           *dg.Channel
	Team2VCh           *dg.Channel
	Match              *Match
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
	list     []*DiscordMatch
	tchMutex sync.RWMutex
	vchMutex sync.RWMutex
	svc      *MatchService
}

func NewDiscordMatchService() *DiscordMatchService {
	return &DiscordMatchService{
		list:     []*DiscordMatch{},
		tchMutex: sync.RWMutex{},
		vchMutex: sync.RWMutex{},
		svc:      NewMatchService(),
	}
}

func (r *DiscordMatchService) Create(tch *dg.Channel, owner *dg.User) (*DiscordMatch, error) {
	r.tchMutex.Lock()
	defer r.tchMutex.Unlock()

	if isContain(tch.ID, du.Channels2IDs(r.getUsingTCh())) {
		return nil, errs.MatchAlreadyStarted
	}

	mt, err := r.svc.Create(&user.User{
		ID:   owner.ID,
		Name: owner.Username,
	})
	if err != nil {
		return nil, err
	}

	dmt := &DiscordMatch{
		Owner:  owner,
		tch:    tch,
		status: StateVCh1Setting,
		Match:  mt,
	}
	r.list = append(r.list, dmt)
	return dmt, nil
}

func (r *DiscordMatchService) Remove(tchID string) error {
	r.tchMutex.Lock()
	defer r.tchMutex.Unlock()
	r.vchMutex.Lock()
	defer r.vchMutex.Unlock()

	for i, mt := range r.list {
		if mt.tch.ID == tchID {
			r.list[i] = r.list[len(r.list)-1]
			r.list[len(r.list)-1] = nil
			r.list = r.list[:len(r.list)-1]
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

func (r *DiscordMatchService) Shuffle(tchID string, vss []*dg.VoiceState) error {
	mt, err := r.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	for _, p := range vss {
		if err != nil {
			return err
		}
		mt.Match.AppendMember(&user.User{
			ID:   p.UserID,
			Name: "",
		})
	}
	return r.svc.Shuffle(mt.Match.Owner)
}

func (r *DiscordMatchService) GetTeam(tchID string) (Team1UserIDs []string, Team2UserIDs []string) {
	mt, err := r.GetMatchByTChID(tchID)
	if err != nil {
		return nil, nil
	}

	f := func(ps []*user.User) []string {
		r := []string{}
		for _, p := range ps {
			r = append(r, p.ID)
		}
		return r
	}
	return f(mt.Match.Team1), f(mt.Match.Team2)
}

// use cache if we need more performance (not map)
func (r *DiscordMatchService) GetMatchByTChID(tchID string) (*DiscordMatch, error) {
	for _, m := range r.list {
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
	for _, m := range r.list {
		if m.tch != nil {
			using = append(using, m.tch)
		}
	}
	return using
}

func (r *DiscordMatchService) getUsingVCh() []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range r.list {
		if m.Team1VCh != nil {
			using = append(using, m.Team1VCh)
		}
		if m.Team2VCh != nil {
			using = append(using, m.Team2VCh)
		}
	}
	return using
}

func isContain(s string, list []string) bool {
	for _, l := range list {
		if s == l {
			return true
		}
	}
	return false
}
