package bot

import (
	"errors"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"team-making-bot/pkg/discord"

	dg "github.com/bwmarrin/discordgo"
	"github.com/golang/glog"
)

var ds *discord.Session
var du *discord.Utility

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
	apiKey string
	mts    *matches
}

func New(apiKey string) *Bot {
	return &Bot{
		apiKey: apiKey,
		mts:    newMatches(),
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

// matchに依存しないように書き換える．manageをifにする
func (b *Bot) onMessageCreate(_ *dg.Session, m *dg.MessageCreate) {
	glog.V(debug).Infof("Channel \"%s\": Get message", *ds.ChannelUnsafe(m.ChannelID))

	if m.Author.ID == ds.State.User.ID {
		glog.V(debug).Infoln("Ignore self message")
		return
	}

	if ds.IsMentionedMessage(m) {
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

	// Handle message if match started on channel
	st, err := b.mts.getMatchStatus(m.ChannelID)
	if errors.Is(err, ers.MatchNotFound) {
		glog.V(info).Infof("Channel \"%s\": can't get match status", *ds.ChannelUnsafe(m.ChannelID), err)
		return
	}
	if err != nil {
		glog.Errorf("Channel \"%s\": can't get match status", *ds.ChannelUnsafe(m.ChannelID), err)
		// Don't send message to user because avoid be troll
		return
	}

	if *st == vCh1Setting || *st == vCh2Setting {
		b.handleVChSettingMessage(m.ChannelID, m.Content, *st)
		return
	}

	if *st == teamPreview {
		// TODO: declare keyword as constant
		if hasKeyword(`shuffle`, m.Content) {
			b.cmdShuffle(m)
			return
		} else if hasKeyword(`go`, m.Content) {
			err := b.mts.movePlayers(m.ChannelID)
			if err != nil {
				glog.Errorf("Channel \"%s\": can't move player %v", err)
				return
			}
			return
		}
	}
	glog.Warningf("Channel \"%s\": Message \"%s\" is not handled", *ds.ChannelUnsafe(m.ChannelID), m.ContentWithMentionsReplaced())
}

func (b *Bot) onMessageReaction(_ *dg.Session, m *dg.MessageReactionAdd) {
	glog.V(debug).Infof("Channel \"%s\": Get reaction", *ds.ChannelUnsafe(m.ChannelID))

	if m.UserID == ds.State.User.ID {
		glog.V(debug).Infoln("Ignore self reaction")
		return
	}
	mt, err := b.mts.getMatchByTChID(m.ChannelID)
	if err != nil {
		glog.V(debug).Infof("Channel \"%s\": Ignore reaction because %s", *ds.ChannelUnsafe(m.ChannelID), err)
		return
	}

	// TODO: Make listening message have callback that handle reaction
	if mt.listeningMessage == nil || mt.listeningMessage.ID != m.MessageID {
		glog.V(debug).Infof("Channel \"%s\": Ignore reaction because message is not listened", mt.tch.Name)
		return
	}

	if mt.status == vCh1Setting {
		if m.Emoji.Name == Stamp["y"] {
			glog.V(debug).Infof("Channel \"%s\": Select yes", mt.tch.Name)
			if mt.recommendedChannel == nil {
				glog.Errorf("Channel \"%s\": Recommend channel is nil", mt.tch.Name)
				ds.ChannelMessageSend(
					m.ChannelID,
					msgs.UnknownError.Format(),
				)
				return
			}

			err = b.mts.setVCh(m.ChannelID, mt.recommendedChannel, "Team1")
			if err == ers.ConflictVCh {
				ds.ChannelMessageSend(
					m.ChannelID,
					msgs.ConflictVCh.Format(mt.recommendedChannel.Name),
				)
				ds.ChannelMessageSend(m.ChannelID, msgs.AskTeam1VCh.Format())
				return
			} else if err != nil {
				glog.Errorf("Channel \"%s\": Cannot set voice channel because %s", mt.tch.Name, err)
				return
			}

			mt.status = vCh2Setting
			mt.recommendedChannel = nil
			mt.listeningMessage = nil

			ds.ChannelMessageSend(
				m.ChannelID,
				msgs.ConfirmTeam1VCh.Format(mt.team1VCh.Name),
			)
			ds.ChannelMessageSend(
				m.ChannelID,
				msgs.AskTeam2VCh.Format(),
			)
			glog.V(info).Infof(
				"Channel \"%s\": Team1 use \"%s\" channel in \"%s\"'s match",
				mt.tch.Name,
				mt.team1VCh.Name,
				mt.tch.Name,
			)
			return
		}
		if m.Emoji.Name == Stamp["n"] {
			mt.listeningMessage = nil
			glog.V(debug).Infof("Channel \"%s\": Select no", mt.tch.Name)
			ds.ChannelMessageSend(
				m.ChannelID,
				msgs.RequestChName.Format(),
			)
			return
		}
	}
	glog.Warningf("Channel \"%s\": Reaction \"%s\" is not handled", *ds.ChannelUnsafe(m.ChannelID), m.Emoji.Name)
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

	tChs := du.FilterChannelsByType(chs, dg.ChannelTypeGuildText)
	if len(tChs) == 0 {
		glog.Warningf("Guild \"%s\": has no text channel", g.Name)
		return
	}
	_, err = ds.ChannelMessageSend(
		tChs[0].ID,
		msgs.Help.Format(),
	)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot send hello message because %s", *ds.ChannelUnsafe(tChs[0].ID), err)
		return
	}
	glog.V(info).Infof("Guild \"%s\": Send hello to \"%s\" channel", g.Name, tChs[0].Name)
}

func (b *Bot) cmdStart(m *dg.MessageCreate) {
	mt, err := b.mts.createMatch(m.ChannelID)
	if errors.Is(err, ers.MatchAlreadyStarted) {
		glog.Warningf("Channel \"%s\": Match already started", *ds.ChannelUnsafe(m.ChannelID))
		ds.ChannelMessageSend(
			m.ChannelID,
			msgs.MatchAlreadyStarted.Format(),
		)
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot handle start command because %s", *ds.ChannelUnsafe(m.ChannelID), err)
		ds.ChannelMessageSend(
			m.ChannelID,
			msgs.UnknownError.Format(),
		)
		return
	}

	err = b.recommendChannel(m.ChannelID)
	if errors.Is(err, ers.NoAvailableVCh) {
		glog.Warningf("Channel \"%s\": Cannot recommend voice channel because %s", mt.tch.Name, err)
		ds.ChannelMessageSend(
			m.ChannelID,
			msgs.NoVChAvailable.Format(),
		)
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot recommend voice channel because %s", mt.tch.Name, err)
		ds.ChannelMessageSend(
			m.ChannelID,
			msgs.UnknownError.Format(),
		)
		return
	}
	glog.V(debug).Infof("Channel \"%s\": Match started", mt.tch.Name)
}

func (b *Bot) cmdExit(m *dg.MessageCreate) {
	err := b.mts.removeMatch(m.ChannelID)
	if errors.Is(err, ers.MatchNotFound) {
		glog.Warningf("Channel \"%s\": Match not found", *ds.ChannelUnsafe(m.ChannelID))
		return
	} else if err != nil {
		glog.Errorf("Channel \"%s\": Cannot handle end command because %s", *ds.ChannelUnsafe(m.ChannelID), err)
		ds.ChannelMessageSend(
			m.ChannelID,
			msgs.UnknownError.Format(),
		)
		return
	}
	ds.ChannelMessageSend(
		m.ChannelID,
		msgs.Exit.Format(),
	)
	glog.V(info).Infof("Channel \"%s\": match has deleted", *ds.ChannelUnsafe(m.ChannelID))
}

func (b *Bot) cmdHelp(m *dg.MessageCreate) {
	_, err := ds.ChannelMessageSend(
		m.ChannelID,
		msgs.Help.Format(),
	)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot handle help command because %s", *ds.ChannelUnsafe(m.ChannelID), err)
		ds.ChannelMessageSend(
			m.ChannelID,
			msgs.UnknownError.Format(),
		)
		return
	}
	glog.V(info).Infof("Channel \"%s\": Help message sent", *ds.ChannelUnsafe(m.ChannelID))
}

func (b *Bot) cmdShuffle(m *dg.MessageCreate) {
	err := b.mts.makeTeam(m.ChannelID)
	if err != nil {
		glog.Errorf("Cannot make team: %v", err)

		return
	}
	err = b.mts.previewTeam(m.ChannelID)
	if err != nil {
		glog.Errorf("Cannot preview team: %v", err)
		ds.ChannelMessageSend(
			m.ChannelID,
			msgs.UnknownError.Format(),
		)
		return
	}
}

func (b *Bot) handleVChSettingMessage(tchID, content string, st Status) {
	tch, err := ds.State.Channel(tchID)
	if err != nil {
		glog.Errorf("Channel %s: Cannot get channel: %v", tchID, err)
		return
	}
	ctx := &struct {
		team       string
		askMsg     *Message
		confirmMsg *Message
	}{
		"Team1",
		msgs.AskTeam1VCh,
		msgs.ConfirmTeam1VCh,
	}
	if st == vCh2Setting {
		ctx.team = "Team2"
		ctx.askMsg = msgs.AskTeam2VCh
		ctx.confirmMsg = msgs.ConfirmTeam2VCh
	}

	chs, err := ds.GetSameGuildChannels(tchID)
	if err != nil {
		glog.Errorf("Channel \"%s\": Cannot get guild channels", tch.Name)
		return
	}
	vChs := du.FilterChannelsByType(chs, dg.ChannelTypeGuildVoice)

	for _, vch := range vChs {
		if hasKeyword(vch.Name, content) {
			err := b.mts.setVCh(tchID, vch, ctx.team)
			if errors.Is(err, ers.ConflictVCh) {
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

			if st == vCh1Setting {
				ds.ChannelMessageSend(tchID, msgs.AskTeam2VCh.Format())
			} else if st == vCh2Setting {
				err := b.mts.makeTeam(tchID)
				if err != nil {
					glog.Errorf("Cannot make team: %v", err)
					return
				}
				err = b.mts.previewTeam(tchID)
				if err != nil {
					glog.Errorf("Cannot preview team: %v", err)
					ds.ChannelMessageSend(
						tchID,
						msgs.UnknownError.Format(),
					)
					return
				}
			}
			return
		}
	}
	glog.V(info).Infof("Channel: \"%s\": Message ignore", *ds.ChannelUnsafe(tchID))
	ds.ChannelMessageSend(tchID, "？")
}

func (b *Bot) recommendChannel(tchID string) error {
	glog.V(trace).Infoln("recommendChannel called")

	tch, err := ds.Channel(tchID)
	if err != nil {
		return err
	}
	chs, err := ds.GuildChannels(tch.GuildID)
	if err != nil {
		return err
	}

	vChs := du.FilterChannelsByType(chs, dg.ChannelTypeGuildVoice)

	availableVChs := []*dg.Channel{}
	usingVchNames := du.Channels2IDs(b.mts.getUsingVCh())
	for _, vCh := range vChs {
		if !isContain(vCh.ID, usingVchNames) {
			availableVChs = append(availableVChs, vCh)
		}
	}

	glog.V(trace).Infof("availableVChs: %+v", availableVChs)
	if len(availableVChs) == 0 {
		return ers.NoAvailableVCh
	}

	g, err := ds.State.Guild(tch.GuildID)
	if err != nil {
		return err
	}

	var vCh *dg.Channel
	if len(g.VoiceStates) > 0 {
		vCh, err = du.GetMostPeopleVCh(availableVChs, g.VoiceStates)
		if err != nil {
			return err
		}
		if vCh == nil {
			vCh = availableVChs[0]
		}
	} else {
		glog.Warningf("Guild %s: No voice states", g.Name)
		vCh = availableVChs[0]
	}

	err = b.mts.setRecommendedChannel(tchID, vCh)
	if err != nil {
		return err
	}

	msg, err := ds.ChannelMessageSend(
		tchID,
		msgs.AskTeam1VChWithRecommend.Format(vCh.Name),
	)
	if err != nil {
		return err
	}

	err = ds.MessageReactionAdd(tchID, msg.ID, Stamp["y"])
	if err != nil {
		glog.Errorln(err)
	}
	err = ds.MessageReactionAdd(tchID, msg.ID, Stamp["n"])
	if err != nil {
		glog.Errorln(err)
	}
	err = b.mts.setListeningMessage(tchID, msg)
	if err != nil {
		glog.Errorf("Channel %s: Can't set listening message: %v", tch.Name, err)
		ds.ChannelMessageSend(
			tch.ID,
			msgs.UnknownError.Format(),
		)
		return err
	}

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
