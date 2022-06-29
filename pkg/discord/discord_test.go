package discord

import (
	"reflect"
	"testing"

	dg "github.com/bwmarrin/discordgo"
)

func TestNew(t *testing.T) {
	type args struct {
		s *dg.Session
	}
	tests := []struct {
		name string
		args args
		want *Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_ChannelUnsafe(t *testing.T) {
	type args struct {
		chID string
	}
	tests := []struct {
		name string
		s    *Session
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ChannelUnsafe(tt.args.chID); got != tt.want {
				t.Errorf("Session.ChannelUnsafe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUtility_Channels2IDs(t *testing.T) {
	type args struct {
		chs []*dg.Channel
	}
	tests := []struct {
		name string
		u    *Utility
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Channels2IDs(tt.args.chs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utility.Channels2IDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUtility_Channels2Names(t *testing.T) {
	type args struct {
		chs []*dg.Channel
	}
	tests := []struct {
		name string
		u    *Utility
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Channels2Names(tt.args.chs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utility.Channels2Names() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUtility_FilterChannelsByType(t *testing.T) {
	type args struct {
		chs []*dg.Channel
		tp  dg.ChannelType
	}
	tests := []struct {
		name string
		u    *Utility
		args args
		want []*dg.Channel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.FilterChannelsByType(tt.args.chs, tt.args.tp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utility.FilterChannelsByType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_IsMentionedMessage(t *testing.T) {
	type args struct {
		m *dg.MessageCreate
	}
	tests := []struct {
		name string
		s    *Session
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsMentionedMessage(tt.args.m); got != tt.want {
				t.Errorf("Session.IsMentionedMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUtility_PackChannelsAndVoiceStates(t *testing.T) {
	type args struct {
		vchs []*dg.Channel
		vss  []*dg.VoiceState
	}
	tests := []struct {
		name    string
		u       *Utility
		args    args
		want    []*ChWithVss
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.PackChannelsAndVoiceStates(tt.args.vchs, tt.args.vss)
			if (err != nil) != tt.wantErr {
				t.Errorf("Utility.PackChannelsAndVoiceStates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utility.PackChannelsAndVoiceStates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetSameGuildChannels(t *testing.T) {
	type args struct {
		chID string
	}
	tests := []struct {
		name    string
		s       Session
		args    args
		want    []*dg.Channel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetSameGuildChannels(tt.args.chID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Session.GetSameGuildChannels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Session.GetSameGuildChannels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUtility_GetMostPeopleVCh(t *testing.T) {
	type args struct {
		vchs []*dg.Channel
		vss  []*dg.VoiceState
	}
	tests := []struct {
		name    string
		u       *Utility
		args    args
		want    *dg.Channel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.GetMostPeopleVCh(tt.args.vchs, tt.args.vss)
			if (err != nil) != tt.wantErr {
				t.Errorf("Utility.GetMostPeopleVCh() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utility.GetMostPeopleVCh() = %v, want %v", got, tt.want)
			}
		})
	}
}
