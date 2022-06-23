package bot

import (
	"errors"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	tm "team-making-bot/internal/team-maker"

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

const (
	info  = 1
	debug = 5
	trace = 6
)

var Stamp = map[string]string{
	"1": "1️⃣",
	"2": "2️⃣",
	"3": "3️⃣",
	"4": "4️⃣",
	"5": "5️⃣",
	"6": "6️⃣",
	"7": "7️⃣",
	"8": "8️⃣",
	"9": "9️⃣",
	"0": "0️⃣",
	"y": "⭕",
	"n": "❌",
}

type Match struct {
	// TODO: dg.channel getter for nil check
	tCh                *dg.Channel
	team1VCh           *dg.Channel
	team2VCh           *dg.Channel
	team1              []*tm.Player
	team2              []*tm.Player
	recommendedChannel *dg.Channel
	listeningMessage   *dg.Message
	status             Status
}

type Bot struct {
	apiKey   string
	session  *dg.Session
	matches  []*Match
	tChMutex sync.RWMutex
	vChMutex sync.RWMutex
	messages *Messages
	errors   *Errors
}

type chWithVss struct {
	ch  *dg.Channel
	vss []*dg.VoiceState
}

func New(apiKey string) *Bot {
	return &Bot{
		apiKey:   apiKey,
		matches:  []*Match{},
		tChMutex: sync.RWMutex{},
		vChMutex: sync.RWMutex{},
		messages: NewMessages(),
		errors:   NewErrors(),
	}
}

func (b *Bot) getUsingTCh() []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range b.matches {
		if m.tCh != nil {
			using = append(using, m.tCh)
		}
	}
	return using
}

func (b *Bot) getUsingVCh() []*dg.Channel {
	using := []*dg.Channel{}
	for _, m := range b.matches {
		if m.team1VCh != nil {
			using = append(using, m.team1VCh)
		}
		if m.team2VCh != nil {
			using = append(using, m.team2VCh)
		}
	}
	return using
}

func (b *Bot) getChannelName(chID string) string {
	ch, err := b.session.State.Channel(chID)
	if err != nil {
		glog.Errorf("Cannot get channel name because %s", err)
		return "Search error"
	}
	return ch.Name
}

func (b *Bot) getMatchByTChID(chID string) (*Match, error) {
	for _, m := range b.matches {
		if m.tCh.ID == chID {
			return m, nil
		}
	}
	return nil, b.errors.MatchNotFound
}

func (b *Bot) channels2IDs(chs []*dg.Channel) []string {
	ids := []string{}
	for _, v := range chs {
		ids = append(ids, v.ID)
	}
	return ids
}

func (b *Bot) channels2Names(ch []*dg.Channel) []string {
	ids := []string{}
	for _, v := range ch {
		ids = append(ids, v.Name)
	}
	return ids
}

func (b *Bot) filterChannelsByType(
	chs []*dg.Channel,
	tp dg.ChannelType,
) []*dg.Channel {
	// TODO: Consider deep copy
	filtered := []*dg.Channel{}
	for _, c := range chs {
		if c.Type == tp {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func (b *Bot) isMentionedMessage(m *dg.MessageCreate) bool {
	for _, user := range m.Mentions {
		if user.ID == b.session.State.User.ID {
			return true
		}
	}
	return false
}

func (b *Bot) packChannelsAndVoiceStates(vChs []*dg.Channel, vss []*dg.VoiceState) ([]*chWithVss, error) {
	glog.V(trace).Infoln("packChannelsAndVoiceStates called")

	targets := []*chWithVss{}
	if len(vChs) == 0 {
		return nil, b.errors.InvalidArgument
	}
	if len(vss) == 0 {
		return nil, b.errors.InvalidArgument
	}

	for _, vc := range vChs {
		targets = append(targets, &chWithVss{vc, []*dg.VoiceState{}})
	}

	for _, vs := range vss {
		for _, tg := range targets {
			if tg.ch.ID == vs.ChannelID {
				tg.vss = append(tg.vss, vs)
			}
		}
	}
	return targets, nil
}

func (b *Bot) getMostPeopleVCh(vChs []*dg.Channel, vss []*dg.VoiceState) *dg.Channel {
	glog.V(trace).Infoln("getMostPeopleVCh called")

	targets, err := b.packChannelsAndVoiceStates(vChs, vss)
	if err != nil {
		glog.Errorf("Cannot pack channels and voice states because %s", err)
		return nil
	}
	max := struct {
		ch    *dg.Channel
		count int
	}{nil, -1}
	for _, tg := range targets {
		if max.count < len(tg.vss) {
			glog.V(trace).Infof("Channel \"%s\" has %d user", tg.ch.Name, len(tg.vss))
			max.ch, max.count = tg.ch, len(tg.vss)
		}
	}
	return max.ch
}

func (b *Bot) recommendChannel(mt *Match) error {
	glog.V(trace).Infoln("recommendChannel called")

	if mt.tCh == nil {
		return b.errors.Unknown
	}

	chs, err := b.session.GuildChannels(mt.tCh.GuildID)
	if err != nil {
		return err
	}

	vChs := b.filterChannelsByType(chs, dg.ChannelTypeGuildVoice)

	availableVChs := []*dg.Channel{}
	for _, vCh := range vChs {
		if !isContain(vCh.ID, b.channels2IDs(b.getUsingVCh())) {
			availableVChs = append(availableVChs, vCh)
		}
	}

	glog.V(trace).Infof("availableVChs: %+v", availableVChs)
	if len(availableVChs) == 0 {
		return b.errors.NoAvailableVCh
	}

	g, err := b.session.State.Guild(mt.tCh.GuildID)
	if err != nil {
		return err
	}

	var vCh *dg.Channel
	if len(g.VoiceStates) > 0 {
		vCh = b.getMostPeopleVCh(availableVChs, g.VoiceStates)
		if vCh == nil {
			vCh = availableVChs[0]
		}
	} else {
		glog.Warningf("Guild %s: No voice states", mt.tCh.Name)
		vCh = availableVChs[0]
	}
	mt.recommendedChannel = vCh

	msg, err := b.session.ChannelMessageSend(
		mt.tCh.ID,
		b.messages.AskTeam1VChWithRecommend.Format(vCh.Name),
	)
	if err != nil {
		return err
	}

	err = b.session.MessageReactionAdd(mt.tCh.ID, msg.ID, Stamp["y"])
	if err != nil {
		glog.Errorln(err)
	}
	err = b.session.MessageReactionAdd(mt.tCh.ID, msg.ID, Stamp["n"])
	if err != nil {
		glog.Errorln(err)
	}
	mt.listeningMessage = msg

	return nil
}

func (b *Bot) setVCh(mt *Match, ch *dg.Channel, team string) error {
	b.vChMutex.Lock()
	defer b.vChMutex.Unlock()
	if isContain(ch.ID, b.channels2IDs(b.getUsingVCh())) {
		return b.errors.ConflictVCh
	}
	if team == "Team1" {
		mt.team1VCh = ch
		return nil
	} else if team == "Team2" {
		mt.team2VCh = ch
		return nil
	}
	return b.errors.InvalidTeam
}

func (b *Bot) cmdStart(m *dg.MessageCreate) {
	// TODO: declare keyword as constant
	if hasKeyword(`start`, m.Content) {
		mt, err := b.createMatch(m)
		if errors.Is(err, b.errors.MatchAlreadyStarted) {
			glog.Warningf("Channel \"%s\": Match already started", b.getChannelName(m.ChannelID))
			b.session.ChannelMessageSend(
				m.ChannelID,
				b.messages.MatchAlreadyStarted.Format(),
			)
			return
		} else if err != nil {
			glog.Errorf("Channel \"%s\": Cannot handle start command because %s", b.getChannelName(m.ChannelID), err)
			b.session.ChannelMessageSend(
				m.ChannelID,
				b.messages.UnknownError.Format(),
			)
			return
		}

		err = b.recommendChannel(mt)
		if errors.Is(err, b.errors.NoAvailableVCh) {
			glog.Warningf("Channel \"%s\": Cannot recommend voice channel because %s", mt.tCh.Name, err)
			b.session.ChannelMessageSend(
				m.ChannelID,
				b.messages.NoVChAvailable.Format(),
			)
			return
		} else if err != nil {
			glog.Errorf("Channel \"%s\": Cannot recommend voice channel because %s", mt.tCh.Name, err)
			b.session.ChannelMessageSend(
				m.ChannelID,
				b.messages.UnknownError.Format(),
			)
			return
		}
		glog.V(debug).Infof("Channel \"%s\": Match started", mt.tCh.Name)
	}
}

func (b *Bot) cmdExit(m *dg.MessageCreate) {
	err := b.exitMatch(m)
	if errors.Is(err, b.errors.MatchNotFound) {
		glog.Warningf("Channel \"%s\": Match not found", b.getChannelName(m.ChannelID))
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot handle end command because %s", b.getChannelName(m.ChannelID), err)
		b.session.ChannelMessageSend(
			m.ChannelID,
			b.messages.UnknownError.Format(),
		)
		return
	}
	b.session.ChannelMessageSend(
		m.ChannelID,
		b.messages.Exit.Format(),
	)
	glog.V(info).Infof("Channel \"%s\": match has deleted", b.getChannelName(m.ChannelID))
}

func (b *Bot) cmdHelp(m *dg.MessageCreate) {
	_, err := b.session.ChannelMessageSend(
		m.ChannelID,
		b.messages.Help.Format(),
	)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot handle help command because %s", b.getChannelName(m.ChannelID), err)
		b.session.ChannelMessageSend(
			m.ChannelID,
			b.messages.UnknownError.Format(),
		)
		return
	}
	glog.V(info).Infof("Channel \"%s\": Help message sent", b.getChannelName(m.ChannelID))
}

func (b *Bot) handleVChSettingMessage(mt *Match, m *dg.MessageCreate) {
	ctx := &struct {
		team       string
		askMsg     *Message
		confirmMsg *Message
		nextStatus Status
	}{
		"Team1",
		b.messages.AskTeam1VCh,
		b.messages.ConfirmTeam1VCh,
		vCh2Setting,
	}
	if mt.status == vCh2Setting {
		ctx.team = "Team2"
		ctx.askMsg = b.messages.AskTeam2VCh
		ctx.confirmMsg = b.messages.ConfirmTeam2VCh
		ctx.nextStatus = teamPreview
	}

	cs, err := b.session.GuildChannels(m.GuildID)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get guild channels", mt.tCh.Name)
		return
	}
	vChs := b.filterChannelsByType(cs, dg.ChannelTypeGuildVoice)

	for _, vCh := range vChs {
		if hasKeyword(vCh.Name, m.Content) {
			err := b.setVCh(mt, vCh, ctx.team)
			if errors.Is(err, b.errors.ConflictVCh) {
				b.session.ChannelMessageSend(
					m.ChannelID,
					b.messages.ConflictVCh.Format(vCh.Name),
				)
				b.session.ChannelMessageSend(
					m.ChannelID,
					ctx.askMsg.Format(),
				)
				glog.Warningf("Channel \"%s\": Conflict voice channel", mt.tCh.Name)
				return
			} else if err != nil {
				glog.Errorf("Channel \"%s\": Cannot set voice channel because %s", mt.tCh.Name, err)
				return
			}

			mt.status = ctx.nextStatus
			mt.recommendedChannel = nil
			glog.V(info).Infof(
				"Channel \"%s\": \"%s\" use \"%s\" channel",
				mt.tCh.Name,
				ctx.team,
				vCh.Name,
			)
			b.session.ChannelMessageSend(
				m.ChannelID,
				ctx.confirmMsg.Format(vCh.Name),
			)
			if mt.status == vCh2Setting {
				b.session.ChannelMessageSend(
					m.ChannelID,
					b.messages.AskTeam2VCh.Format(),
				)
			} else if mt.status == teamPreview {
				b.makeTeam(mt)
				b.previewTeam(mt)
			}
			return
		}
	}
	glog.V(info).Infof("Channel: \"%s\": Message ignore", b.getChannelName(m.ChannelID))
	glog.V(info).Infof("\"%s\" sent trash message! lol", m.Author.Username)
	b.session.ChannelMessageSend(
		m.ChannelID,
		"？",
	)
}

func (b *Bot) Run() {
	session, err := dg.New("Bot " + b.apiKey)
	if err != nil {
		glog.Errorln("Cannot create bot session")
		return
	}
	b.session = session

	b.session.StateEnabled = true
	b.session.AddHandler(b.onMessageCreate)
	b.session.AddHandler(b.onMessageReaction)
	b.session.AddHandler(func(s *dg.Session, r *dg.Ready) {
		for _, g := range s.State.Guilds {
			b.onEnable(g.ID)
		}
	})

	err = b.session.Open()
	if err != nil {
		glog.Fatalln("Cannot open session")
		return
	}
	glog.V(info).Infoln("Session started")

	defer func() {
		err := b.session.Close()
		if err != nil {
			glog.Fatal("Cannot close session")
		}
		glog.V(info).Infoln("Session closed")
		os.Exit(0)
	}()

	stopBot := make(chan os.Signal, 1)

	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	glog.V(info).Infoln("Bot started")
	<-stopBot
}

func (b *Bot) onMessageCreate(s *dg.Session, m *dg.MessageCreate) {
	glog.V(debug).Infof("Channel \"%s\": Get message", b.getChannelName(m.ChannelID))

	if m.Author.ID == s.State.User.ID {
		glog.V(debug).Infoln("Ignore self message")
		return
	}

	if b.isMentionedMessage(m) {
		// TODO: declare keyword as constant
		if hasKeyword(`start`, m.Content) {
			b.cmdStart(m)
			return
		} else if hasKeyword(`end`, m.Content) {
			b.cmdExit(m)
			return
		} else if hasKeyword(`help`, m.Content) {
			b.cmdHelp(m)
			return
		} else if hasKeyword(`reset`, m.Content) {
			b.cmdStart(m)
			b.cmdExit(m)
			return
		}
	}

	// Handle message if setting voice channel state
	mt, err := b.getMatchByTChID(m.ChannelID)
	if errors.Is(err, b.errors.MatchNotFound) {
		glog.V(info).Infof("Channel \"%s\": Ignore message because %s", b.getChannelName(m.ChannelID), err)
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Ignore message because %s", b.getChannelName(m.ChannelID), err)
		return
	}

	if mt.status == vCh1Setting || mt.status == vCh2Setting {
		b.handleVChSettingMessage(mt, m)
		return
	}

	if mt.status == teamPreview {
		// TODO: declare keyword as constant
		if hasKeyword(`shuffle`, m.Content) {
			b.makeTeam(mt)
			b.previewTeam(mt)
			return
		} else if hasKeyword(`go`, m.Content) {
			b.movePlayers(mt)
			return
		}
	}
	glog.Warningf("Channel \"%s\": Message \"%s\" is not handled", b.getChannelName(m.ChannelID), m.ContentWithMentionsReplaced())
}

func (b *Bot) onMessageReaction(s *dg.Session, m *dg.MessageReactionAdd) {
	glog.V(debug).Infof("Channel \"%s\": Get reaction", b.getChannelName(m.ChannelID))

	if m.UserID == s.State.User.ID {
		glog.V(debug).Infoln("Ignore self reaction")
		return
	}
	mt, err := b.getMatchByTChID(m.ChannelID)
	if err != nil {
		glog.V(debug).Infof("Channel \"%s\": Ignore reaction because %s", b.getChannelName(m.ChannelID), err)
		return
	}

	// TODO: Make listening message have callback that handle reaction
	if mt.listeningMessage == nil || mt.listeningMessage.ID != m.MessageID {
		glog.V(debug).Infof("Channel \"%s\": Ignore reaction because message is not listened", mt.tCh.Name)
		return
	}

	if mt.status == vCh1Setting {
		if m.Emoji.Name == Stamp["y"] {
			glog.V(debug).Infof("Channel \"%s\": Select yes", mt.tCh.Name)
			if mt.recommendedChannel == nil {
				glog.Errorf("Channel \"%s\": Recommend channel is nil", mt.tCh.Name)
				b.session.ChannelMessageSend(
					m.ChannelID,
					b.messages.UnknownError.Format(mt.recommendedChannel.Name),
				)
				return
			}

			err = b.setVCh(mt, mt.recommendedChannel, "Team1")
			if err == b.errors.ConflictVCh {
				b.session.ChannelMessageSend(
					m.ChannelID,
					b.messages.ConflictVCh.Format(mt.recommendedChannel.Name),
				)
				b.session.ChannelMessageSend(
					m.ChannelID,
					b.messages.AskTeam1VCh.Format(),
				)
				return
			} else if err != nil {
				glog.Errorf("Channel \"%s\": Cannot set voice channel because %s", mt.tCh.Name, err)
				return
			}

			mt.status = vCh2Setting
			mt.recommendedChannel = nil
			mt.listeningMessage = nil

			b.session.ChannelMessageSend(
				m.ChannelID,
				b.messages.ConfirmTeam1VCh.Format(mt.team1VCh.Name),
			)
			b.session.ChannelMessageSend(
				m.ChannelID,
				b.messages.AskTeam2VCh.Format(),
			)
			glog.V(info).Infof(
				"Channel \"%s\": Team1 use \"%s\" channel in \"%s\"'s match",
				mt.tCh.Name,
				mt.team1VCh.Name,
				mt.tCh.Name,
			)
			return
		}
		if m.Emoji.Name == Stamp["n"] {
			mt.listeningMessage = nil
			glog.V(debug).Infof("Channel \"%s\": Select no", mt.tCh.Name)
			b.session.ChannelMessageSend(
				m.ChannelID,
				b.messages.RequestChName.Format(),
			)
			return
		}
	}
	glog.Warningf("Channel \"%s\": Reaction \"%s\" is not handled", b.getChannelName(m.ChannelID), m.Emoji.Name)
}

func (b *Bot) onEnable(g *dg.Guild) {
	chs, err := b.session.GuildChannels(g.ID)
	if err != nil {
		glog.Errorf("GuildID %s: Cannot get channels", g.ID)
		return
	}
	if len(chs) == 0 {
		glog.V(debug).Infof("Guild \"%s\": has no channel", g.Name)
		return
	}

	tChs := b.filterChannelsByType(chs, dg.ChannelTypeGuildText)
	if len(tChs) == 0 {
		glog.Warningf("Guild \"%s\": has no text channel", g.Name)
		return
	}
	_, err = b.session.ChannelMessageSend(
		tChs[0].ID,
		b.messages.Help.Format(),
	)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot send hello message because %s", b.getChannelName(tChs[0].ID), err)
		return
	}
	glog.V(info).Infof("Guild \"%s\": Send hello to \"%s\" channel", g.Name, tChs[0].Name)
}

func (b *Bot) createMatch(m *dg.MessageCreate) (*Match, error) {
	ch, err := b.session.State.Channel(m.ChannelID)
	if err != nil {
		return nil, err
	}

	b.tChMutex.Lock()
	defer b.tChMutex.Unlock()

	if isContain(m.ChannelID, b.channels2IDs(b.getUsingTCh())) {
		return nil, b.errors.MatchAlreadyStarted
	}

	mt := &Match{
		tCh:    ch,
		status: vCh1Setting,
	}
	b.matches = append(b.matches, mt)
	return mt, nil
}

func (b *Bot) exitMatch(m *dg.MessageCreate) error {
	b.tChMutex.Lock()
	defer b.tChMutex.Unlock()
	b.vChMutex.Lock()
	defer b.vChMutex.Unlock()

	for i, mt := range b.matches {
		if mt.tCh.ID == m.ChannelID {
			b.matches[i] = b.matches[len(b.matches)-1]
			b.matches[len(b.matches)-1] = nil
			b.matches = b.matches[:len(b.matches)-1]
			return nil
		}
	}
	return b.errors.MatchNotFound
}

func (b *Bot) makeTeam(mt *Match) {
	g, err := b.session.State.Guild(mt.tCh.GuildID)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get guild", mt.tCh.Name)
		return
	}

	chAndVss, err := b.packChannelsAndVoiceStates([]*dg.Channel{mt.team1VCh, mt.team2VCh}, g.VoiceStates)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot pack channels and voice states because %s", mt.tCh.Name, err)
		b.session.ChannelMessageSend(
			mt.tCh.ID,
			b.messages.UnknownError.Format(),
		)
		return
	}
	if len(chAndVss) != 2 {
		glog.Errorf("Channel \"%s\": Cannot get voice state", mt.tCh.Name)
		b.session.ChannelMessageSend(
			mt.tCh.ID,
			b.messages.UnknownError.Format(),
		)
		return
	}

	b.session.ChannelMessageSend(
		mt.tCh.ID,
		b.messages.MakeTeam.Format(mt.team1VCh.Name, mt.team2VCh.Name),
	)

	players := []*tm.Player{}
	for _, cv := range chAndVss {
		for _, p := range cv.vss {
			usr, err := b.session.User(p.UserID)
			if err != nil {
				glog.Errorf("User: \"%s\" not found", p.UserID)
				return
			}
			players = append(players, &tm.Player{DiscordId: p.UserID, Name: usr.Username})
		}
	}
	rtm, err := tm.NewRandomTeamMaker()
	if err != nil {
		println(err)
	}
	mt.team1, mt.team2 = rtm.MakeTeam(players)
}

func (b *Bot) previewTeam(mt *Match) {
	f := func(ps []*tm.Player) []string {
		r := []string{}
		for _, p := range ps {
			r = append(r, p.Name)
		}
		return r
	}
	names1, names2 := f(mt.team1), f(mt.team2)
	b.session.ChannelMessageSend(
		mt.tCh.ID,
		b.messages.ConfirmTeam.Format(mt.team1VCh.Name, names1, mt.team2VCh.Name, names2),
	)
}

func (b *Bot) movePlayers(mt *Match) {
	for _, p := range mt.team1 {
		err := b.session.GuildMemberMove(mt.tCh.GuildID, p.DiscordId, &mt.team1VCh.ID)
		if err != nil {
			glog.Errorf("Channel \"%s\": Cannot move member because %s", mt.tCh.Name, err)
			b.session.ChannelMessageSend(
				mt.tCh.ID,
				b.messages.UnknownError.Format(),
			)
			return
		}
	}
	for _, p := range mt.team2 {
		err := b.session.GuildMemberMove(mt.tCh.GuildID, p.DiscordId, &mt.team2VCh.ID)
		if err != nil {
			glog.Errorf("Channel \"%s\": Cannot move member because %s", mt.tCh.Name, err)
			b.session.ChannelMessageSend(
				mt.tCh.ID,
				b.messages.UnknownError.Format(),
			)
			return
		}
	}
}

func isContain(s string, list []string) bool {
	for _, l := range list {
		if s == l {
			return true
		}
	}
	return false
}

func hasKeyword(keyword, target string) bool {
	r, err := regexp.Compile(keyword)
	if err != nil {
		glog.Errorln("Cannot compile regex")
		return false
	}
	return r.MatchString(target)
}
