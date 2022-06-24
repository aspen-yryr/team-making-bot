package discord

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
)

type Session struct {
	*dg.Session
}

type ChWithVss struct {
	Ch  *dg.Channel
	Vss []*dg.VoiceState
}

type Utility struct {
}

func New(s *dg.Session) *Session {
	return &Session{s}
}

// No error handling
func (s *Session) ChannelUnsafe(chID string) *string {
	ch, err := s.State.Channel(chID)
	if err != nil {
		return nil
	}
	return &ch.Name
}

func (u *Utility) Channels2IDs(chs []*dg.Channel) []string {
	ids := []string{}
	for _, v := range chs {
		ids = append(ids, v.ID)
	}
	return ids
}

func (u *Utility) Channels2Names(chs []*dg.Channel) []string {
	names := []string{}
	for _, v := range chs {
		names = append(names, v.Name)
	}
	return names
}

func (u *Utility) FilterChannelsByType(
	chs []*dg.Channel,
	tp dg.ChannelType,
) []*dg.Channel {
	filtered := []*dg.Channel{}
	for _, c := range chs {
		if c.Type == tp {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func (s *Session) IsMentionedMessage(m *dg.MessageCreate) bool {
	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID {
			return true
		}
	}
	return false
}

// TODO: add log to reference
func (u *Utility) PackChannelsAndVoiceStates(vchs []*dg.Channel, vss []*dg.VoiceState) ([]*ChWithVss, error) {
	pk := []*ChWithVss{}
	if len(vchs) == 0 {
		return nil, fmt.Errorf("no voice channels in arg")
	}
	if len(vss) == 0 {
		return nil, fmt.Errorf("no voice states in arg")
	}

	for _, vc := range vchs {
		pk = append(pk, &ChWithVss{vc, []*dg.VoiceState{}})
	}

	for _, vs := range vss {
		for _, tg := range pk {
			if tg.Ch.ID == vs.ChannelID {
				tg.Vss = append(tg.Vss, vs)
			}
		}
	}
	return pk, nil
}

func (s Session) GetSameGuildChannels(chID string) ([]*dg.Channel, error) {
	ch, err := s.Channel(chID)
	if err != nil {
		return nil, fmt.Errorf("can't get channel: %v", err)
	}
	chs, err := s.GuildChannels(ch.GuildID)
	if err != nil {
		return nil, fmt.Errorf("can't get guild channels: %v", err)
	}
	return chs, nil
}

func (u *Utility) GetMostPeopleVCh(vChs []*dg.Channel, vss []*dg.VoiceState) (*dg.Channel, error) {
	targets, err := u.PackChannelsAndVoiceStates(vChs, vss)
	if err != nil {
		return nil, fmt.Errorf("can't pack channels and voice states: %v", err)
	}
	max := struct {
		ch    *dg.Channel
		count int
	}{nil, -1}
	for _, tg := range targets {
		if max.count < len(tg.Vss) {
			max.ch, max.count = tg.Ch, len(tg.Vss)
		}
	}
	return max.ch, nil
}
