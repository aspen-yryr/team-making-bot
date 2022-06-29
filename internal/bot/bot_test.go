package bot

import (
	"reflect"
	"team-making-bot/internal/match"
	"testing"

	dg "github.com/bwmarrin/discordgo"
)

func TestNew(t *testing.T) {
	type args struct {
		apiKey string
		greet  bool
	}
	tests := []struct {
		name string
		args args
		want *Bot
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.apiKey, tt.args.greet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBot_Run(t *testing.T) {
	tests := []struct {
		name string
		b    *Bot
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Run()
		})
	}
}

func TestBot_onMessageCreate(t *testing.T) {
	type args struct {
		in0 *dg.Session
		m   *dg.MessageCreate
	}
	tests := []struct {
		name string
		b    *Bot
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.onMessageCreate(tt.args.in0, tt.args.m)
		})
	}
}

func TestBot_onMessageReaction(t *testing.T) {
	type args struct {
		in0 *dg.Session
		m   *dg.MessageReactionAdd
	}
	tests := []struct {
		name string
		b    *Bot
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.onMessageReaction(tt.args.in0, tt.args.m)
		})
	}
}

func TestBot_onEnable(t *testing.T) {
	type args struct {
		g *dg.Guild
	}
	tests := []struct {
		name string
		b    *Bot
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.onEnable(tt.args.g)
		})
	}
}

func TestBot_cmdStart(t *testing.T) {
	type args struct {
		m *dg.MessageCreate
	}
	tests := []struct {
		name string
		b    *Bot
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.cmdStart(tt.args.m)
		})
	}
}

func TestBot_cmdExit(t *testing.T) {
	type args struct {
		m *dg.MessageCreate
	}
	tests := []struct {
		name string
		b    *Bot
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.cmdExit(tt.args.m)
		})
	}
}

func TestBot_cmdHelp(t *testing.T) {
	type args struct {
		m *dg.MessageCreate
	}
	tests := []struct {
		name string
		b    *Bot
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.cmdHelp(tt.args.m)
		})
	}
}

func TestBot_cmdShuffle(t *testing.T) {
	type args struct {
		m *dg.MessageCreate
	}
	tests := []struct {
		name string
		b    *Bot
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.cmdShuffle(tt.args.m)
		})
	}
}

func TestBot_handleVChSettingMessage(t *testing.T) {
	type args struct {
		tchID   string
		content string
		st      match.Status
	}
	tests := []struct {
		name string
		b    *Bot
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.handleVChSettingMessage(tt.args.tchID, tt.args.content, tt.args.st)
		})
	}
}

func TestBot_recommendChannel(t *testing.T) {
	type args struct {
		tchID string
	}
	tests := []struct {
		name    string
		b       *Bot
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.recommendChannel(tt.args.tchID); (err != nil) != tt.wantErr {
				t.Errorf("Bot.recommendChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBot_movePlayers(t *testing.T) {
	type args struct {
		tchID string
	}
	tests := []struct {
		name    string
		b       *Bot
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.movePlayers(tt.args.tchID); (err != nil) != tt.wantErr {
				t.Errorf("Bot.movePlayers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBot_previewTeam(t *testing.T) {
	type args struct {
		tchID string
	}
	tests := []struct {
		name    string
		b       *Bot
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.previewTeam(tt.args.tchID); (err != nil) != tt.wantErr {
				t.Errorf("Bot.previewTeam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_hasKeyword(t *testing.T) {
	type args struct {
		keyword string
		target  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasKeyword(tt.args.keyword, tt.args.target); got != tt.want {
				t.Errorf("hasKeyword() = %v, want %v", got, tt.want)
			}
		})
	}
}
