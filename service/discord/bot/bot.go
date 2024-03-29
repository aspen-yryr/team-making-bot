package bot

import (
	"errors"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/aspen-yryr/team-making-bot/internal/constants"
	"github.com/aspen-yryr/team-making-bot/pkg/discord"
	"github.com/aspen-yryr/team-making-bot/service/discord/match"

	dg "github.com/bwmarrin/discordgo"
	"github.com/golang/glog"
)

var ds *discord.Session
var msgs = constants.Msgs
var errs = constants.Errs

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

type Bot struct {
	apiKey   string
	matchSvc *match.DiscordMatchService
	greet    bool
}

func New(apiKey string, greet bool) *Bot {
	return &Bot{
		apiKey:   apiKey,
		matchSvc: match.NewDiscordMatchService(),
		greet:    greet,
	}
}

func (b *Bot) Run() {
	session, err := dg.New("Bot " + b.apiKey)
	if err != nil {
		glog.Errorln("Cannot create bot session")
		return
	}
	ds = discord.New(session)

	ds.StateEnabled = true
	ds.AddHandler(b.onMessageCreate)
	ds.AddHandler(b.onMessageReaction)
	ds.AddHandler(func(_ *dg.Session, r *dg.Ready) {
		for _, g := range ds.State.Guilds {
			b.onEnable(g)
		}
	})

	err = ds.Open()
	if err != nil {
		glog.Fatalln("Cannot open session")
		return
	}
	glog.V(info).Infoln("Session started")

	defer func() {
		err := ds.Close()
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

func (b *Bot) onMessageCreate(_ *dg.Session, msg *dg.MessageCreate) {
	glog.V(debug).Infof("Channel \"%s\": Get message", ds.ChannelUnsafe(msg.ChannelID))

	if msg.Author.ID == ds.State.User.ID {
		glog.V(debug).Infoln("Ignore self message")
		return
	}

	if ds.IsMentionedMessage(msg) {
		// TODO: declare keyword as constant
		if hasKeyword(`start`, msg.Content) {
			b.cmdStart(msg)
			return
		} else if hasKeyword(`end`, msg.Content) {
			b.cmdExit(msg)
			return
		} else if hasKeyword(`help`, msg.Content) {
			b.cmdHelp(msg)
			return
		} else if hasKeyword(`reset`, msg.Content) {
			b.cmdStart(msg)
			b.cmdExit(msg)
			return
		}
	}

	mt, err := b.matchSvc.GetMatchByTChID(msg.ChannelID)
	if errors.Is(err, errs.MatchNotFound) {
		glog.V(info).Infof("Channel \"%s\": can't get match status: %v", ds.ChannelUnsafe(msg.ChannelID), err)
		return
	}
	if err != nil {
		glog.Errorf("Channel \"%s\": can't get match status: %v", ds.ChannelUnsafe(msg.ChannelID), err)
		// Don't send message to user because avoid be troll
		return
	}
	status := mt.Status

	if status == match.StateVCh1Setting || status == match.StateVCh2Setting {
		b.handleVChSettingMessage(msg.ChannelID, msg.Content, status)
		return
	}

	if status == match.StateTeamPreview {
		// TODO: declare keyword as constant
		if hasKeyword(`shuffle`, msg.Content) {
			b.cmdShuffle(msg)
			return
		} else if hasKeyword(`go`, msg.Content) {
			err := b.movePlayers(msg.ChannelID)
			if err != nil {
				glog.Errorf("Channel \"%s\": can't move player %v", err)
				return
			}
			return
		}
	}
	glog.Warningf("Channel \"%s\": Message \"%s\" is not handled", ds.ChannelUnsafe(msg.ChannelID), msg.ContentWithMentionsReplaced())
}

func (b *Bot) onMessageReaction(_ *dg.Session, msg *dg.MessageReactionAdd) {
	glog.V(debug).Infof("Channel \"%s\": Get reaction", ds.ChannelUnsafe(msg.ChannelID))

	if msg.UserID == ds.State.User.ID {
		glog.V(debug).Infoln("Ignore self reaction")
		return
	}

	mt, err := b.matchSvc.GetMatchByTChID(msg.ChannelID)
	if err != nil {
		glog.V(debug).Infof("Channel \"%s\": Ignore reaction because %s", ds.ChannelUnsafe(msg.ChannelID), err)
		return
	}

	status := mt.Status

	if err != nil {
		glog.Errorf("Channel \"%s\": can't get match status", ds.ChannelUnsafe(msg.ChannelID), err)
		// Don't send message to user because avoid be troll
		return
	}
	if status != match.StateVCh1Setting {
		return
	}

	// TODO: Make listening message have callback that handle reaction
	if !b.matchSvc.IsListeningMessage(msg.ChannelID, msg.MessageID) {
		glog.V(debug).Infof("Channel \"%s\": Ignore reaction because message is not listened", ds.ChannelUnsafe(msg.ChannelID))
		return
	}

	if msg.Emoji.Name == Stamp["y"] {
		glog.V(debug).Infof("Channel \"%s\": Select yes", ds.ChannelUnsafe(msg.ChannelID))
		vch, err := mt.GetRecommendedChannel()
		if err != nil {
			glog.Errorf("Channel \"%s\": can't get recommended channel", ds.ChannelUnsafe(msg.ChannelID))
			ds.ChannelMessageSend(
				msg.ChannelID,
				msgs.UnknownError.Format(),
			)
			return
		}

		err = b.matchSvc.SetVCh(msg.ChannelID, vch, "Team1")
		if err == errs.ConflictVCh {
			ds.ChannelMessageSend(
				msg.ChannelID,
				msgs.ConflictVCh.Format(vch.Name),
			)
			ds.ChannelMessageSend(msg.ChannelID, msgs.AskTeam1VCh.Format())
			return
		} else if err != nil {
			glog.Errorf("Channel \"%s\": Cannot set voice channel because %s", ds.ChannelUnsafe(msg.ChannelID), err)
			return
		}

		glog.V(info).Infof(
			"Channel \"%s\": Team1 use \"%s\" channel",
			ds.ChannelUnsafe(msg.ChannelID),
			vch.Name,
		)
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.ConfirmTeam1VCh.Format(vch.Name),
		)
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.AskTeam2VCh.Format(),
		)
		return
	}
	if msg.Emoji.Name == Stamp["n"] {
		glog.V(debug).Infof("Channel \"%s\": Select no", ds.ChannelUnsafe(msg.ChannelID))
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.RequestChName.Format(),
		)
		return
	}
	glog.Warningf("Channel \"%s\": Reaction \"%s\" is not handled", ds.ChannelUnsafe(msg.ChannelID), msg.Emoji.Name)
}
func (b *Bot) onEnable(g *dg.Guild) {
	chs, err := ds.GuildChannels(g.ID)
	if err != nil {
		glog.Errorf("GuildID %s: Cannot get channels", g.ID)
		return
	}
	if len(chs) == 0 {
		glog.V(debug).Infof("Guild \"%s\": has no channel", g.Name)
		return
	}

	tchs := ds.FilterChannelsByType(chs, dg.ChannelTypeGuildText)
	if len(tchs) == 0 {
		glog.Warningf("Guild \"%s\": has no text channel", g.Name)
		return
	}
	if b.greet {
		_, err = ds.ChannelMessageSend(
			tchs[0].ID,
			msgs.Help.Format(),
		)
		if err != nil {
			glog.Errorf("Channel \"%s\": Cannot send hello message because %s", ds.ChannelUnsafe(tchs[0].ID), err)
			return
		}
		glog.V(info).Infof("Guild \"%s\": Send hello to \"%s\" channel", g.Name, tchs[0].Name)
	} else {
		glog.V(info).Infof("Guild \"%s\": Bot activate", g.Name)
	}
}

func (b *Bot) cmdStart(msg *dg.MessageCreate) {
	tch, err := ds.Channel(msg.ChannelID)
	if err != nil {
		glog.Errorf("can't get channel: %v", err)
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.UnknownError.Format(),
		)
		return
	}

	// TODO: canCreate method
	vch, err := b.getOwnerVch(msg.ChannelID, msg.Author.ID)
	if errors.Is(err, errs.OwnerNotInVchs) {
		glog.Errorf("Channel \"%s\": Cannot get owner voice channel because %s", ds.ChannelUnsafe(msg.ChannelID), err)
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.OwnerNotInVchs.Format(),
		)
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get owner voice channel because %s", ds.ChannelUnsafe(msg.ChannelID), err)
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.UnknownError.Format(),
		)
	}

	_, err = b.matchSvc.Create(tch, msg.Author)
	if errors.Is(err, errs.MatchAlreadyStarted) {
		glog.Warningf("Channel \"%s\": Match already started", ds.ChannelUnsafe(msg.ChannelID))
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.MatchAlreadyStarted.Format(),
		)
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot handle start command because %s", ds.ChannelUnsafe(msg.ChannelID), err)
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.UnknownError.Format(),
		)
		b.matchSvc.Remove(tch.ID)
		return
	}

	chs, err := ds.GuildChannels(tch.GuildID)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get guild channels because %s", ds.ChannelUnsafe(msg.ChannelID), err)
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.UnknownError.Format(),
		)
		b.matchSvc.Remove(tch.ID)
		return
	}
	availableVChs := b.matchSvc.FilterAvailableVCh(chs)
	if len(availableVChs) == 0 {
		glog.Warningf("Channel \"%s\": No voice channel available", ds.ChannelUnsafe(msg.ChannelID))
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.NoVChAvailable.Format(),
		)
		b.matchSvc.Remove(tch.ID)
		return
	}

	err = b.recommendChannel(msg.ChannelID, vch.ID)
	if errors.Is(err, errs.NoAvailableVCh) {
		glog.Warningf("Channel \"%s\": Cannot recommend voice channel because %s", ds.ChannelUnsafe(msg.ChannelID), err)
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.NoVChAvailable.Format(),
		)
		b.matchSvc.Remove(tch.ID)
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot recommend voice channel because %s", ds.ChannelUnsafe(msg.ChannelID), err)
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.UnknownError.Format(),
		)
		b.matchSvc.Remove(tch.ID)
		return
	}
	glog.V(debug).Infof("Channel \"%s\": Match started", ds.ChannelUnsafe(msg.ChannelID))
}

func (b *Bot) cmdExit(msg *dg.MessageCreate) {
	err := b.matchSvc.Remove(msg.ChannelID)
	if errors.Is(err, errs.MatchNotFound) {
		glog.Warningf("Channel \"%s\": Match not found", ds.ChannelUnsafe(msg.ChannelID))
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot handle end command because %s", ds.ChannelUnsafe(msg.ChannelID), err)
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.UnknownError.Format(),
		)
		return
	}
	ds.ChannelMessageSend(
		msg.ChannelID,
		msgs.Exit.Format(),
	)
	glog.V(info).Infof("Channel \"%s\": match has deleted", ds.ChannelUnsafe(msg.ChannelID))
}

func (b *Bot) cmdHelp(msg *dg.MessageCreate) {
	_, err := ds.ChannelMessageSend(
		msg.ChannelID,
		msgs.Help.Format(),
	)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot handle help command because %s", ds.ChannelUnsafe(msg.ChannelID), err)
		ds.ChannelMessageSend(
			msg.ChannelID,
			msgs.UnknownError.Format(),
		)
		return
	}
	glog.V(info).Infof("Channel \"%s\": Help message sent", ds.ChannelUnsafe(msg.ChannelID))
}

func (b *Bot) cmdShuffle(msg *dg.MessageCreate) {
	err := b._shuffle(msg.GuildID, msg.ChannelID)
	if err != nil {
		glog.Errorf("Cannot shuffle team: %v", err)
		return
	}
}

func (b *Bot) handleVChSettingMessage(tchID, content string, status match.Status) {
	tch, err := ds.State.Channel(tchID)
	if err != nil {
		glog.Errorf("Channel %s: Cannot get channel: %v", tchID, err)
		return
	}

	ctx := &struct {
		team       string
		askMsg     *constants.Message
		confirmMsg *constants.Message
	}{
		"Team1",
		msgs.AskTeam1VCh,
		msgs.ConfirmTeam1VCh,
	}
	if status == match.StateVCh2Setting {
		ctx.team = "Team2"
		ctx.askMsg = msgs.AskTeam2VCh
		ctx.confirmMsg = msgs.ConfirmTeam2VCh
	}

	chs, err := ds.GetSameGuildChannels(tchID)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get guild channels", tch.Name)
		return
	}
	vchs := ds.FilterChannelsByType(chs, dg.ChannelTypeGuildVoice)

	for _, vch := range vchs {
		if hasKeyword(vch.Name, content) {
			err := b.matchSvc.SetVCh(tchID, vch, ctx.team)
			if errors.Is(err, errs.ConflictVCh) {
				ds.ChannelMessageSend(tchID, msgs.ConflictVCh.Format(vch.Name))
				ds.ChannelMessageSend(tchID, ctx.askMsg.Format())
				glog.Warningf("Channel \"%s\": Conflict voice channel", tch.Name)
				return
			} else if err != nil {
				glog.Errorf("Channel \"%s\": Cannot set voice channel because %s", tch.Name, err)
				return
			}

			glog.V(info).Infof(
				"Channel \"%s\": \"%s\" use \"%s\" channel",
				tch.Name,
				ctx.team,
				vch.Name,
			)
			ds.ChannelMessageSend(tchID, ctx.confirmMsg.Format(vch.Name))

			if status == match.StateVCh1Setting {
				ds.ChannelMessageSend(tchID, msgs.AskTeam2VCh.Format())
			} else if status == match.StateVCh2Setting {
				g, err := ds.State.Guild(tch.GuildID)
				if err != nil {
					glog.Errorf("can't get guild: $v", err)
					return
				}
				if len(g.VoiceStates) > 0 {
					err = b._shuffle(g.ID, tchID)
					if err != nil {
						glog.Errorf("can't shuffle team: %v", err)
						return
					}
				}

			}
			return
		}
	}
	glog.V(info).Infof("Channel: \"%s\": Message ignore", ds.ChannelUnsafe(tchID))
	ds.ChannelMessageSend(tchID, "？")
}

func (b *Bot) getOwnerVch(tchID string, owner_id string) (*dg.Channel, error) {
	glog.V(trace).Infoln("getOwnerVch called")
	tch, err := ds.Channel(tchID)
	if err != nil {
		return nil, err
	}
	g, err := ds.State.Guild(tch.GuildID)
	if err != nil {
		return nil, err
	}

	var vch *dg.Channel = nil
	for _, vs := range g.VoiceStates {
		if vs.UserID == owner_id {
			vch, err = ds.Channel(vs.ChannelID)
			if err != nil {
				return nil, err
			}
		}
	}
	if vch == nil {
		return nil, errs.OwnerNotInVchs
	}

	return vch, nil
}

func (b *Bot) recommendChannel(tchID, vchID string) error {
	glog.V(trace).Infoln("recommendChannel called")

	vch, err := ds.Channel(vchID)
	if err != nil {
		return err
	}

	msg, err := ds.ChannelMessageSend(
		tchID,
		msgs.AskTeam1VChWithRecommend.Format(vch.Name),
	)
	if err != nil {
		return err
	}

	mt, err := b.matchSvc.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}
	mt.RecommendedChannel = vch

	err = ds.MessageReactionAdd(tchID, msg.ID, Stamp["y"])
	if err != nil {
		glog.Errorln(err)
	}
	err = ds.MessageReactionAdd(tchID, msg.ID, Stamp["n"])
	if err != nil {
		glog.Errorln(err)
	}
	err = b.matchSvc.SetListeningMessage(tchID, msg)
	if err != nil {
		glog.Errorf("Channel %s: Can't set listening message: %v", ds.ChannelUnsafe(tchID), err)
		ds.ChannelMessageSend(
			tchID,
			msgs.UnknownError.Format(),
		)
		return err
	}

	return nil
}

func (b *Bot) _shuffle(gID, tchID string) error {
	g, err := ds.State.Guild(gID)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get guild", ds.ChannelUnsafe(tchID))
		return err
	}

	mt, err := b.matchSvc.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	var chWithVss []*discord.ChWithVss
	chWithVss, err = ds.PackChannelsAndVoiceStates([]*dg.Channel{mt.Team1VCh, mt.Team2VCh}, g.VoiceStates)
	if err != nil {
		glog.Errorf("Cannot pack voice states: %v", err)
		return err
	}

	players := []*dg.User{}
	for _, cv := range chWithVss {
		for _, p := range cv.Vss {
			u, err := ds.User(p.UserID)
			if err != nil {
				glog.Errorf("Cannot get discord user: %v", err)
				return err
			}
			players = append(players, u)
		}
	}

	err = b.matchSvc.AppendMembers(tchID, players)
	if err != nil {
		glog.Errorf("Cannot append members: %v", err)
		return err
	}

	err = b.matchSvc.Shuffle(tchID)
	if err != nil {
		glog.Errorf("Cannot shuffle team: %v", err)
		return err
	}
	err = b.previewTeam(tchID)
	if err != nil {
		glog.Errorf("Cannot preview team: %v", err)
		ds.ChannelMessageSend(
			tchID,
			msgs.UnknownError.Format(),
		)
		return err
	}
	return nil
}

func (b *Bot) movePlayers(tchID string) error {
	mt, err := b.matchSvc.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	team1, team2, err := b.matchSvc.GetTeam(tchID)
	if err != nil {
		return err
	}

	for _, p := range team1 {
		err := ds.GuildMemberMove(mt.GetGuildId(), p, &mt.Team1VCh.ID)
		if err != nil {
			glog.Errorf("Channel \"%s\": Cannot move member because %s", ds.ChannelUnsafe(tchID), err)
			ds.ChannelMessageSend(tchID, msgs.UnknownError.Format())
			return err
		}
	}
	for _, p := range team2 {
		err := ds.GuildMemberMove(mt.GetGuildId(), p, &mt.Team2VCh.ID)
		if err != nil {
			glog.Errorf("Channel \"%s\": Cannot move member because %s", ds.ChannelUnsafe(tchID), err)
			ds.ChannelMessageSend(tchID, msgs.UnknownError.Format())
			return err
		}
	}
	return nil
}

func (b *Bot) previewTeam(tchID string) error {
	ids1, ids2, err := b.matchSvc.GetTeam(tchID)
	if err != nil {
		return err
	}
	mt, err := b.matchSvc.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}
	f := func(ids []string) []string {
		r := []string{}
		for _, id := range ids {
			u, _ := ds.User(id)
			r = append(r, u.Username)
		}
		return r
	}
	ds.ChannelMessageSend(
		tchID,
		msgs.ConfirmTeam.Format(mt.Team1VCh.Name, f(ids1), mt.Team2VCh.Name, f(ids2)),
	)
	return nil
}

func hasKeyword(keyword, target string) bool {
	r, err := regexp.Compile(keyword)
	if err != nil {
		glog.Errorln("Cannot compile regex")
		return false
	}
	return r.MatchString(target)
}
