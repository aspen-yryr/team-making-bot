package bot

import (
	"errors"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/aspen-yryr/team-making-bot/internal/configs"
	"github.com/aspen-yryr/team-making-bot/internal/constants"
	"github.com/aspen-yryr/team-making-bot/internal/database/migration"
	"github.com/aspen-yryr/team-making-bot/internal/user"
	"github.com/aspen-yryr/team-making-bot/pkg/dg_wrap"
	"github.com/aspen-yryr/team-making-bot/service/discord/match"
	dg "github.com/bwmarrin/discordgo"

	"github.com/golang/glog"
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

//DiscordとMatchの橋渡し
type Bot struct {
	conf       configs.BotConfig
	discordSvc dg_wrap.DiscordSvc
	matchSvc   *match.DiscordMatchService
	userRepo   user.UserRepository
	errs       constants.Errors
	msgs       constants.Messages
}

func NewBot(
	conf *configs.BotConfig,
	ds *dg_wrap.DiscordSvc,
	ms *match.DiscordMatchService,
	ur user.UserRepository,
	errs *constants.Errors,
	msg *constants.Messages,
	m migration.Migrator,
) *Bot {
	err := m.Run()
	if err != nil {
		panic("migration  failed")
	}
	return &Bot{
		conf:       *conf,
		discordSvc: *ds,
		matchSvc:   ms,
		userRepo:   ur,
		errs:       *errs,
		msgs:       *msg,
	}
}

func (b *Bot) Run() {
	b.discordSvc.AddHandler(b.onMessageCreate)
	b.discordSvc.AddHandler(b.onMessageReaction)
	// b.discordSvc.AddHandler(func(_ dg.Session, r dg.Ready) {
	// 	for _, g := range b.discordSvc.State.Guilds {
	// 		b.onEnable(g)
	// 	}
	// })

	err := b.discordSvc.Open()
	if err != nil {
		glog.Fatalln("Cannot open session")
		return
	}
	glog.V(info).Infoln("Session started")

	defer func() {
		err := b.discordSvc.Close()
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
	glog.V(debug).Infof("Channel \"%s\": Get message", b.discordSvc.ChannelUnsafe(msg.ChannelID))

	if msg.Author.ID == b.discordSvc.State.User.ID {
		glog.V(debug).Infoln("Ignore self message")
		return
	}

	if b.discordSvc.IsMentionedMessage(msg) {
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
	if errors.Is(err, b.errs.MatchNotFound) {
		glog.V(info).Infof("Channel \"%s\": can't get match status: %v", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
		return
	}
	if err != nil {
		glog.Errorf("Channel \"%s\": can't get match status: %v", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
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
	glog.Warningf("Channel \"%s\": Message \"%s\" is not handled", b.discordSvc.ChannelUnsafe(msg.ChannelID), msg.ContentWithMentionsReplaced())
}

func (b *Bot) onMessageReaction(_ *dg.Session, msg *dg.MessageReactionAdd) {
	glog.V(debug).Infof("Channel \"%s\": Get reaction", b.discordSvc.ChannelUnsafe(msg.ChannelID))

	if msg.UserID == b.discordSvc.State.User.ID {
		glog.V(debug).Infoln("Ignore self reaction")
		return
	}

	mt, err := b.matchSvc.GetMatchByTChID(msg.ChannelID)
	if err != nil {
		glog.V(debug).Infof("Channel \"%s\": Ignore reaction because %s", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
		return
	}

	status := mt.Status

	if err != nil {
		glog.Errorf("Channel \"%s\": can't get match status", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
		// Don't send message to user because avoid be troll
		return
	}
	if status != match.StateVCh1Setting {
		return
	}

	// TODO: Make listening message have callback that handle reaction
	if !b.matchSvc.IsListeningMessage(msg.ChannelID, msg.MessageID) {
		glog.V(debug).Infof("Channel \"%s\": Ignore reaction because message is not listened", b.discordSvc.ChannelUnsafe(msg.ChannelID))
		return
	}

	if msg.Emoji.Name == Stamp["y"] {
		glog.V(debug).Infof("Channel \"%s\": Select yes", b.discordSvc.ChannelUnsafe(msg.ChannelID))
		vch, err := mt.GetRecommendedChannel()
		if err != nil {
			glog.Errorf("Channel \"%s\": can't get recommended channel", b.discordSvc.ChannelUnsafe(msg.ChannelID))
			b.discordSvc.ChannelMessageSend(
				msg.ChannelID,
				b.msgs.UnknownError.Format(),
			)
			return
		}

		err = b.matchSvc.SetVCh(msg.ChannelID, vch, "Team1")
		if err == b.errs.ConflictVCh {
			b.discordSvc.ChannelMessageSend(
				msg.ChannelID,
				b.msgs.ConflictVCh.Format(vch.Name),
			)
			b.discordSvc.ChannelMessageSend(msg.ChannelID, b.msgs.AskTeam1VCh.Format())
			return
		} else if err != nil {
			glog.Errorf("Channel \"%s\": Cannot set voice channel because %s", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
			return
		}

		glog.V(info).Infof(
			"Channel \"%s\": Team1 use \"%s\" channel",
			b.discordSvc.ChannelUnsafe(msg.ChannelID),
			vch.Name,
		)
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.ConfirmTeam1VCh.Format(vch.Name),
		)
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.AskTeam2VCh.Format(),
		)
		return
	}
	if msg.Emoji.Name == Stamp["n"] {
		glog.V(debug).Infof("Channel \"%s\": Select no", b.discordSvc.ChannelUnsafe(msg.ChannelID))
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.RequestChName.Format(),
		)
		return
	}
	glog.Warningf("Channel \"%s\": Reaction \"%s\" is not handled", b.discordSvc.ChannelUnsafe(msg.ChannelID), msg.Emoji.Name)
}

func (b *Bot) onEnable(g *dg.Guild) {
	chs, err := b.discordSvc.GuildChannels(g.ID)
	if err != nil {
		glog.Errorf("GuildID %s: Cannot get channels", g.ID)
		return
	}
	if len(chs) == 0 {
		glog.V(debug).Infof("Guild \"%s\": has no channel", g.Name)
		return
	}

	tchs := b.discordSvc.FilterChannelsByType(chs, dg.ChannelTypeGuildText)
	if len(tchs) == 0 {
		glog.Warningf("Guild \"%s\": has no text channel", g.Name)
		return
	}
	if b.conf.Greet {
		_, err = b.discordSvc.ChannelMessageSend(
			tchs[0].ID,
			b.msgs.Help.Format(),
		)
		if err != nil {
			glog.Errorf("Channel \"%s\": Cannot send hello message because %s", b.discordSvc.ChannelUnsafe(tchs[0].ID), err)
			return
		}
		glog.V(info).Infof("Guild \"%s\": Send hello to \"%s\" channel", g.Name, tchs[0].Name)
	} else {
		glog.V(info).Infof("Guild \"%s\": Bot activate", g.Name)
	}
}

func (b *Bot) cmdStart(msg *dg.MessageCreate) {
	tch, err := b.discordSvc.Channel(msg.ChannelID)
	if err != nil {
		glog.Errorf("can't get channel: %v", err)
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.UnknownError.Format(),
		)
		return
	}

	// TODO: canCreate method
	vch, err := b.getOwnerVch(msg.ChannelID, msg.Author.ID)
	if errors.Is(err, b.errs.OwnerNotInVchs) {
		glog.Errorf("Channel \"%s\": Cannot get owner voice channel because %s", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.OwnerNotInVchs.Format(),
		)
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get owner voice channel because %s", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.UnknownError.Format(),
		)
	}

	_, err = b.matchSvc.Create(tch, msg.Author)
	if errors.Is(err, b.errs.MatchAlreadyStarted) {
		glog.Warningf("Channel \"%s\": Match already started", b.discordSvc.ChannelUnsafe(msg.ChannelID))
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.MatchAlreadyStarted.Format(),
		)
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot handle start command because %s", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.UnknownError.Format(),
		)
		b.matchSvc.Remove(tch.ID)
		return
	}

	chs, err := b.discordSvc.GuildChannels(tch.GuildID)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get guild channels because %s", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.UnknownError.Format(),
		)
		b.matchSvc.Remove(tch.ID)
		return
	}
	availableVChs := b.matchSvc.FilterAvailableVCh(chs)
	if len(availableVChs) == 0 {
		glog.Warningf("Channel \"%s\": No voice channel available", b.discordSvc.ChannelUnsafe(msg.ChannelID))
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.NoVChAvailable.Format(),
		)
		b.matchSvc.Remove(tch.ID)
		return
	}

	err = b.recommendChannel(msg.ChannelID, vch.ID)
	if errors.Is(err, b.errs.NoAvailableVCh) {
		glog.Warningf("Channel \"%s\": Cannot recommend voice channel because %s", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.NoVChAvailable.Format(),
		)
		b.matchSvc.Remove(tch.ID)
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot recommend voice channel because %s", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.UnknownError.Format(),
		)
		b.matchSvc.Remove(tch.ID)
		return
	}
	glog.V(debug).Infof("Channel \"%s\": Match started", b.discordSvc.ChannelUnsafe(msg.ChannelID))
}

func (b *Bot) cmdExit(msg *dg.MessageCreate) {
	err := b.matchSvc.Remove(msg.ChannelID)
	if errors.Is(err, b.errs.MatchNotFound) {
		glog.Warningf("Channel \"%s\": Match not found", b.discordSvc.ChannelUnsafe(msg.ChannelID))
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot handle end command because %s", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.UnknownError.Format(),
		)
		return
	}
	b.discordSvc.ChannelMessageSend(
		msg.ChannelID,
		b.msgs.Exit.Format(),
	)
	glog.V(info).Infof("Channel \"%s\": match has deleted", b.discordSvc.ChannelUnsafe(msg.ChannelID))
}

func (b *Bot) cmdHelp(msg *dg.MessageCreate) {
	_, err := b.discordSvc.ChannelMessageSend(
		msg.ChannelID,
		b.msgs.Help.Format(),
	)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot handle help command because %s", b.discordSvc.ChannelUnsafe(msg.ChannelID), err)
		b.discordSvc.ChannelMessageSend(
			msg.ChannelID,
			b.msgs.UnknownError.Format(),
		)
		return
	}
	glog.V(info).Infof("Channel \"%s\": Help message sent", b.discordSvc.ChannelUnsafe(msg.ChannelID))
}

func (b *Bot) cmdShuffle(msg *dg.MessageCreate) {
	err := b._shuffle(msg.GuildID, msg.ChannelID)
	if err != nil {
		glog.Errorf("Cannot shuffle team: %v", err)
		return
	}
}

func (b *Bot) handleVChSettingMessage(tchID, content string, status match.Status) {
	tch, err := b.discordSvc.State.Channel(tchID)
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
		b.msgs.AskTeam1VCh,
		b.msgs.ConfirmTeam1VCh,
	}
	if status == match.StateVCh2Setting {
		ctx.team = "Team2"
		ctx.askMsg = b.msgs.AskTeam2VCh
		ctx.confirmMsg = b.msgs.ConfirmTeam2VCh
	}

	chs, err := b.discordSvc.GetSameGuildChannels(tchID)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get guild channels", tch.Name)
		return
	}
	vchs := b.discordSvc.FilterChannelsByType(chs, dg.ChannelTypeGuildVoice)

	for _, vch := range vchs {
		if hasKeyword(vch.Name, content) {
			err := b.matchSvc.SetVCh(tchID, vch, ctx.team)
			if errors.Is(err, b.errs.ConflictVCh) {
				b.discordSvc.ChannelMessageSend(tchID, b.msgs.ConflictVCh.Format(vch.Name))
				b.discordSvc.ChannelMessageSend(tchID, ctx.askMsg.Format())
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
			b.discordSvc.ChannelMessageSend(tchID, ctx.confirmMsg.Format(vch.Name))

			if status == match.StateVCh1Setting {
				b.discordSvc.ChannelMessageSend(tchID, b.msgs.AskTeam2VCh.Format())
			} else if status == match.StateVCh2Setting {
				g, err := b.discordSvc.State.Guild(tch.GuildID)
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
	glog.V(info).Infof("Channel: \"%s\": Message ignore", b.discordSvc.ChannelUnsafe(tchID))
	b.discordSvc.ChannelMessageSend(tchID, "？")
}

func (b *Bot) getOwnerVch(tchID string, owner_id string) (*dg.Channel, error) {
	glog.V(trace).Infoln("getOwnerVch called")
	tch, err := b.discordSvc.Channel(tchID)
	if err != nil {
		return nil, err
	}
	g, err := b.discordSvc.State.Guild(tch.GuildID)
	if err != nil {
		return nil, err
	}

	var vch *dg.Channel = nil
	for _, vs := range g.VoiceStates {
		if vs.UserID == owner_id {
			vch, err = b.discordSvc.Channel(vs.ChannelID)
			if err != nil {
				return nil, err
			}
		}
	}
	if vch == nil {
		return nil, b.errs.OwnerNotInVchs
	}

	return vch, nil
}

func (b *Bot) recommendChannel(tchID, vchID string) error {
	glog.V(trace).Infoln("recommendChannel called")

	vch, err := b.discordSvc.Channel(vchID)
	if err != nil {
		return err
	}

	msg, err := b.discordSvc.ChannelMessageSend(
		tchID,
		b.msgs.AskTeam1VChWithRecommend.Format(vch.Name),
	)
	if err != nil {
		return err
	}

	mt, err := b.matchSvc.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}
	mt.RecommendedChannel = vch

	err = b.discordSvc.MessageReactionAdd(tchID, msg.ID, Stamp["y"])
	if err != nil {
		glog.Errorln(err)
	}
	err = b.discordSvc.MessageReactionAdd(tchID, msg.ID, Stamp["n"])
	if err != nil {
		glog.Errorln(err)
	}
	err = b.matchSvc.SetListeningMessage(tchID, msg)
	if err != nil {
		glog.Errorf("Channel %s: Can't set listening message: %v", b.discordSvc.ChannelUnsafe(tchID), err)
		b.discordSvc.ChannelMessageSend(
			tchID,
			b.msgs.UnknownError.Format(),
		)
		return err
	}
	return nil
}

func (b *Bot) _shuffle(gID, tchID string) error {
	g, err := b.discordSvc.State.Guild(gID)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get guild", b.discordSvc.ChannelUnsafe(tchID))
		return err
	}

	mt, err := b.matchSvc.GetMatchByTChID(tchID)
	if err != nil {
		return err
	}

	var chWithVss []*dg_wrap.ChWithVss
	chWithVss, err = b.discordSvc.PackChannelsAndVoiceStates([]*dg.Channel{mt.Team1VCh, mt.Team2VCh}, g.VoiceStates)
	if err != nil {
		glog.Errorf("Cannot pack voice states: %v", err)
		return err
	}

	players := []*dg.User{}
	for _, cv := range chWithVss {
		for _, p := range cv.Vss {
			u, err := b.discordSvc.User(p.UserID)
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
		b.discordSvc.ChannelMessageSend(
			tchID,
			b.msgs.UnknownError.Format(),
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
		err := b.discordSvc.GuildMemberMove(mt.GetGuildId(), p, &mt.Team1VCh.ID)
		if err != nil {
			glog.Errorf("Channel \"%s\": Cannot move member because %s", b.discordSvc.ChannelUnsafe(tchID), err)
			b.discordSvc.ChannelMessageSend(tchID, b.msgs.UnknownError.Format())
			return err
		}
	}
	for _, p := range team2 {
		err := b.discordSvc.GuildMemberMove(mt.GetGuildId(), p, &mt.Team2VCh.ID)
		if err != nil {
			glog.Errorf("Channel \"%s\": Cannot move member because %s", b.discordSvc.ChannelUnsafe(tchID), err)
			b.discordSvc.ChannelMessageSend(tchID, b.msgs.UnknownError.Format())
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
			u, _ := b.discordSvc.User(id)
			r = append(r, u.Username)
		}
		return r
	}
	b.discordSvc.ChannelMessageSend(
		tchID,
		b.msgs.ConfirmTeam.Format(mt.Team1VCh.Name, f(ids1), mt.Team2VCh.Name, f(ids2)),
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
