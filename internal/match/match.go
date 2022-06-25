package match

import (
	"sync"
	"team-making-bot/internal/constants"
	"team-making-bot/pkg/discord"
	tm "team-making-bot/pkg/team-maker"

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

type Match struct {
	// TODO: dg.channel getter for nil check
	tch                *dg.Channel
	Team1VCh           *dg.Channel
	Team2VCh           *dg.Channel
	team1              []*tm.Player
	team2              []*tm.Player
	recommendedChannel *dg.Channel
	listeningMessage   *dg.Message
	status             Status
}

//TODO: More strictly state management or more get method to nil check

func (m *Match) GetRecommendedChannel() (*dg.Channel, error) {
	if m.recommendedChannel == nil {
		return nil, errs.Unknown
	}
	return m.recommendedChannel, nil
}

func (m *Match) GetShuffledIds() ([]string, []string) {
	f := func(ps []*tm.Player) []string {
		r := []string{}
		for _, p := range ps {
			r = append(r, p.DiscordId)
		}
		return r
	}
	return f(m.team1), f(m.team2)
}

func (m *Match) GetGuildId() string {
	return m.tch.GuildID
}

type Manager struct {
	list     []*Match
	tchMutex sync.RWMutex
	vchMutex sync.RWMutex
}

func NewMatches() *Manager {
	return &Manager{
		list:     []*Match{},
		tchMutex: sync.RWMutex{},
		vchMutex: sync.RWMutex{},
	}
}

func (m *Manager) CreateMatch(tch *dg.Channel) (*Match, error) {
	m.tchMutex.Lock()
	defer m.tchMutex.Unlock()

	if isContain(tch.ID, du.Channels2IDs(m.getUsingTCh(false))) {
		return nil, errs.MatchAlreadyStarted
	}

	mt := &Match{
		tch:    tch,
		status: StateVCh1Setting,
	}
	m.list = append(m.list, mt)
	return mt, nil
}

func (m *Manager) RemoveMatch(tchID string) error {
	m.tchMutex.Lock()
	defer m.tchMutex.Unlock()
	m.vchMutex.Lock()
	defer m.vchMutex.Unlock()

	for i, mt := range m.list {
		if mt.tch.ID == tchID {
			m.list[i] = m.list[len(m.list)-1]
			m.list[len(m.list)-1] = nil
			m.list = m.list[:len(m.list)-1]
			return nil
		}
	}
	return errs.MatchNotFound
}

func (m *Manager) FilterAvailableVCh(chs []*dg.Channel) []*dg.Channel {
	vchs := du.FilterChannelsByType(chs, dg.ChannelTypeGuildVoice)
	if len(vchs) == 0 {
		return []*dg.Channel{}
	}

	availableVChs := []*dg.Channel{}
	for _, vch := range vchs {
		if !isContain(vch.ID, du.Channels2IDs(m.getUsingVCh(true))) {
			availableVChs = append(availableVChs, vch)
		}
	}
	return availableVChs
}

func (m *Manager) SetVCh(tchID string, vch *dg.Channel, team string) error {
	mt, err := m.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	m.vchMutex.Lock()
	defer m.vchMutex.Unlock()

	if isContain(vch.ID, du.Channels2IDs(m.getUsingVCh(false))) {
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

func (m *Manager) ShuffleTeam(tchID string, vss []*dg.VoiceState) error {
	mt, err := m.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	var chWithVss []*discord.ChWithVss
	chWithVss, _ = du.PackChannelsAndVoiceStates([]*dg.Channel{mt.Team1VCh, mt.Team2VCh}, vss)

	players := []*tm.Player{}
	for _, cv := range chWithVss {
		for _, p := range cv.Vss {
			if err != nil {
				return err
			}
			players = append(players, &tm.Player{DiscordId: p.UserID})
		}
	}
	rtm, err := tm.NewRandomTeamMaker()
	if err != nil {
		return err
	}
	mt.team1, mt.team2 = rtm.MakeTeam(players)
	return nil
}

func (m *Manager) GetTeam(tchID string) (Team1UserIDs []string, Team2UserIDs []string) {
	mt, err := m.GetMatchByTChID(tchID)
	if err != nil {
		return nil, nil
	}

	f := func(ps []*tm.Player) []string {
		r := []string{}
		for _, p := range ps {
			r = append(r, p.DiscordId)
		}
		return r
	}
	return f(mt.team1), f(mt.team2)
}

// use cache if we need more performance (not map)
func (m *Manager) GetMatchByTChID(tchID string) (*Match, error) {
	for _, m := range m.list {
		if m.tch.ID == tchID {
			return m, nil
		}
	}
	return nil, errs.MatchNotFound
}

func (m *Manager) GetMatchStatus(tchID string) (*Status, error) {
	mt, err := m.GetMatchByTChID(tchID)
	if err != nil {
		return nil, err
	}
	return &mt.status, nil
}

func (m *Manager) SetRecommendedChannel(tchID string, vch *dg.Channel) error {
	mt, err := m.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}
	mt.recommendedChannel = vch
	return nil
}

func (m *Manager) SetListeningMessage(tchID string, msg *dg.Message) error {
	mt, err := m.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}
	mt.listeningMessage = msg
	return nil
}

func (m *Manager) IsListeningMessage(tchID, msgID string) bool {
	mt, err := m.GetMatchByTChID(tchID)
	if err != nil {
		return false
	}
	if mt.listeningMessage.ID == msgID {
		return true
	}
	return false
}

// use cache if we need more performance(not map)
func (m *Manager) getUsingTCh(lock bool) []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range m.list {
		if m.tch != nil {
			using = append(using, m.tch)
		}
	}
	return using
}

func (m *Manager) getUsingVCh(lock bool) []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range m.list {
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
