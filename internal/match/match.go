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

func (m *Match) StatusIs(r Status) bool {
	return m.status == r
}

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
	tChMutex sync.RWMutex
	vChMutex sync.RWMutex
}

func NewMatches() *Manager {
	return &Manager{
		list:     []*Match{},
		tChMutex: sync.RWMutex{},
		vChMutex: sync.RWMutex{},
	}
}

func (m *Manager) CreateMatch(tch *dg.Channel) (*Match, error) {
	m.tChMutex.Lock()
	defer m.tChMutex.Unlock()

	if isContain(tch.ID, du.Channels2IDs(m.getUsingTCh())) {
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
	m.tChMutex.Lock()
	defer m.tChMutex.Unlock()
	m.vChMutex.Lock()
	defer m.vChMutex.Unlock()

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

func (m *Manager) ShuffleTeam(tchID string, vss []*dg.VoiceState) error {
	mt, err := m.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	var chWithVss []*discord.ChWithVss
	chWithVss, _ = du.PackChannelsAndVoiceStates([]*dg.Channel{mt.Team1VCh, mt.Team2VCh}, vss)
	// if err != nil {
	// 	return err
	// }

	// ds.ChannelMessageSend(
	// 	mt.tch.ID,
	// 	msgs.MakeTeam.Format(mt.team1VCh.Name, mt.team2VCh.Name),
	// )

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

func (m *Manager) SetVCh(tchID string, vch *dg.Channel, team string) error {
	mt, err := m.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	m.vChMutex.Lock()
	defer m.vChMutex.Unlock()

	if isContain(vch.ID, du.Channels2IDs(m.getUsingVCh())) {
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

func (m *Manager) getUsingTCh() []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range m.list {
		if m.tch != nil {
			using = append(using, m.tch)
		}
	}
	return using
}

func (m *Manager) getUsingVCh() []*dg.Channel {
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

func (m *Manager) FilterAvailableVCh(chs []*dg.Channel) []*dg.Channel {
	vChs := du.FilterChannelsByType(chs, dg.ChannelTypeGuildVoice)
	if len(vChs) == 0 {
		return []*dg.Channel{}
	}

	availableVChs := []*dg.Channel{}
	for _, vCh := range vChs {
		if !isContain(vCh.ID, du.Channels2IDs(m.getUsingVCh())) {
			availableVChs = append(availableVChs, vCh)
		}
	}
	return availableVChs
}

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

func isContain(s string, list []string) bool {
	for _, l := range list {
		if s == l {
			return true
		}
	}
	return false
}
