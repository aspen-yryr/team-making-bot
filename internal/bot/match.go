package bot

import (
	"sync"
	"team-making-bot/pkg/discord"
	tm "team-making-bot/pkg/team-maker"

	dg "github.com/bwmarrin/discordgo"
	"github.com/golang/glog"
)

type (
	Status int
)

const (
	vCh1Setting Status = iota
	vCh2Setting
	teamPreview
)

type match struct {
	// TODO: dg.channel getter for nil check
	tch                *dg.Channel
	team1VCh           *dg.Channel
	team2VCh           *dg.Channel
	team1              []*tm.Player
	team2              []*tm.Player
	recommendedChannel *dg.Channel
	listeningMessage   *dg.Message
	status             Status
}

type matches struct {
	list     []*match
	tChMutex sync.RWMutex
	vChMutex sync.RWMutex
}

func newMatches() *matches {
	return &matches{
		list:     []*match{},
		tChMutex: sync.RWMutex{},
		vChMutex: sync.RWMutex{},
	}
}

func (mts *matches) createMatch(tchID string) (*match, error) {
	tch, err := ds.State.Channel(tchID)
	if err != nil {
		return nil, err
	}
	mts.tChMutex.Lock()
	defer mts.tChMutex.Unlock()

	if isContain(tchID, du.Channels2IDs(mts.getUsingTCh())) {
		return nil, ers.MatchAlreadyStarted
	}

	mt := &match{
		tch:    tch,
		status: vCh1Setting,
	}
	mts.list = append(mts.list, mt)
	return mt, nil
}

func (mts *matches) removeMatch(tchID string) error {
	mts.tChMutex.Lock()
	defer mts.tChMutex.Unlock()
	mts.vChMutex.Lock()
	defer mts.vChMutex.Unlock()

	for i, mt := range mts.list {
		if mt.tch.ID == tchID {
			mts.list[i] = mts.list[len(mts.list)-1]
			mts.list[len(mts.list)-1] = nil
			mts.list = mts.list[:len(mts.list)-1]
			return nil
		}
	}
	return ers.MatchNotFound
}

func (mts *matches) makeTeam(tchID string) error {
	mt, err := mts.getMatchByTChID(tchID)
	if err != nil {
		glog.Errorf("Channel \"%s\": can't get match: %v", mt.tch.Name, err)
		return err
	}
	g, err := ds.State.Guild(mt.tch.GuildID)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get guild", mt.tch.Name)
		return err
	}
	var chWithVss []*discord.ChWithVss
	if len(g.VoiceStates) > 0 {
		chWithVss, err = du.PackChannelsAndVoiceStates([]*dg.Channel{mt.team1VCh, mt.team2VCh}, g.VoiceStates)
		if err != nil {
			glog.Errorf("Channel \"%s\": Cannot pack channels and voice states because %s", mt.tch.Name, err)
			return err
		}
		if len(chWithVss) != 2 {
			glog.Errorf("Channel \"%s\": Cannot get voice state", mt.tch.Name)
			return err
		}
	} else {
		glog.Warningf("Guild \"%s\": No voice states", g.Name)
	}
	ds.ChannelMessageSend(
		mt.tch.ID,
		msgs.MakeTeam.Format(mt.team1VCh.Name, mt.team2VCh.Name),
	)

	players := []*tm.Player{}
	for _, cv := range chWithVss {
		for _, p := range cv.Vss {
			usr, err := ds.User(p.UserID)
			if err != nil {
				glog.Errorf("User: \"%s\" not found: %v", p.UserID, err)
				return err
			}
			players = append(players, &tm.Player{DiscordId: p.UserID, Name: usr.Username})
		}
	}
	rtm, err := tm.NewRandomTeamMaker()
	if err != nil {
		println(err)
	}
	mt.team1, mt.team2 = rtm.MakeTeam(players)
	return nil
}

func (mts *matches) previewTeam(tchID string) error {
	mt, err := mts.getMatchByTChID(tchID)
	if err != nil {
		return err
	}

	f := func(ps []*tm.Player) []string {
		r := []string{}
		for _, p := range ps {
			r = append(r, p.Name)
		}
		return r
	}
	names1, names2 := f(mt.team1), f(mt.team2)
	ds.ChannelMessageSend(
		mt.tch.ID,
		msgs.ConfirmTeam.Format(mt.team1VCh.Name, names1, mt.team2VCh.Name, names2),
	)
	return nil
}

func (mts *matches) movePlayers(tchID string) error {
	mt, err := mts.getMatchByTChID(tchID)
	if err != nil {
		return err
	}

	for _, p := range mt.team1 {
		err := ds.GuildMemberMove(mt.tch.GuildID, p.DiscordId, &mt.team1VCh.ID)
		if err != nil {
			glog.Errorf("Channel \"%s\": Cannot move member because %s", mt.tch.Name, err)
			ds.ChannelMessageSend(
				mt.tch.ID,
				msgs.UnknownError.Format(),
			)
			return err
		}
	}
	for _, p := range mt.team2 {
		err := ds.GuildMemberMove(mt.tch.GuildID, p.DiscordId, &mt.team2VCh.ID)
		if err != nil {
			glog.Errorf("Channel \"%s\": Cannot move member because %s", mt.tch.Name, err)
			ds.ChannelMessageSend(
				mt.tch.ID,
				msgs.UnknownError.Format(),
			)
			return err
		}
	}
	return nil
}

func (mts *matches) setVCh(tchID string, vch *dg.Channel, team string) error {
	mt, err := mts.getMatchByTChID(tchID)
	if err != nil {
		return err
	}

	mts.vChMutex.Lock()
	defer mts.vChMutex.Unlock()

	if isContain(vch.ID, du.Channels2IDs(mts.getUsingVCh())) {
		return ers.ConflictVCh
	}

	if team == "Team1" {
		mt.team1VCh = vch
		mt.status = vCh2Setting
		mt.recommendedChannel = nil
		return nil
	} else if team == "Team2" {
		mt.team2VCh = vch
		mt.status = teamPreview
		mt.recommendedChannel = nil
		return nil
	}
	return ers.InvalidTeam
}

func (mts *matches) getUsingTCh() []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range mts.list {
		if m.tch != nil {
			using = append(using, m.tch)
		}
	}
	return using
}

func (mts *matches) getUsingVCh() []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range mts.list {
		if m.team1VCh != nil {
			using = append(using, m.team1VCh)
		}
		if m.team2VCh != nil {
			using = append(using, m.team2VCh)
		}
	}
	return using
}

func (mts *matches) getMatchByTChID(tchID string) (*match, error) {
	for _, m := range mts.list {
		if m.tch.ID == tchID {
			return m, nil
		}
	}
	return nil, ers.MatchNotFound
}

func (mts *matches) getMatchStatus(tchID string) (*Status, error) {
	mt, err := mts.getMatchByTChID(tchID)
	if err != nil {
		return nil, err
	}
	return &mt.status, nil
}

func (mts *matches) setRecommendedChannel(tchID string, vch *dg.Channel) error {
	mt, err := mts.getMatchByTChID(tchID)
	if err != nil {
		return err
	}
	mt.recommendedChannel = vch
	return nil
}

func (mts *matches) setListeningMessage(tchID string, msg *dg.Message) error {
	mt, err := mts.getMatchByTChID(tchID)
	if err != nil {
		return err
	}
	mt.listeningMessage = msg
	return nil
}

func isContain(s string, list []string) bool {
	for _, l := range list {
		if s == l {
			return true
		}
	}
	return false
}
