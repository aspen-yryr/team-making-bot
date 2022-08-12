package match

import (
	"reflect"
	"testing"

	dg "github.com/bwmarrin/discordgo"
)

func TestMatch_GetRecommendedChannel(t *testing.T) {
	tests := []struct {
		name    string
		m       *Match
		want    *dg.Channel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetRecommendedChannel()
			if (err != nil) != tt.wantErr {
				t.Errorf("Match.GetRecommendedChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Match.GetRecommendedChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatch_GetShuffledIds(t *testing.T) {
	tests := []struct {
		name  string
		m     *Match
		want  []string
		want1 []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.m.GetShuffledIds()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Match.GetShuffledIds() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Match.GetShuffledIds() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMatch_GetGuildId(t *testing.T) {
	tests := []struct {
		name string
		m    *Match
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetGuildId(); got != tt.want {
				t.Errorf("Match.GetGuildId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMatches(t *testing.T) {
	tests := []struct {
		name string
		want *Manager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMatches(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_CreateMatch(t *testing.T) {
	type args struct {
		tch *dg.Channel
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		want    *Match
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.CreateMatch(tt.args.tch)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.CreateMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.CreateMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_RemoveMatch(t *testing.T) {
	type args struct {
		tchID string
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.RemoveMatch(tt.args.tchID); (err != nil) != tt.wantErr {
				t.Errorf("Manager.RemoveMatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_FilterAvailableVCh(t *testing.T) {
	type args struct {
		chs []*dg.Channel
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want []*dg.Channel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.FilterAvailableVCh(tt.args.chs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.FilterAvailableVCh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_SetVCh(t *testing.T) {
	type args struct {
		tchID string
		vch   *dg.Channel
		team  string
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SetVCh(tt.args.tchID, tt.args.vch, tt.args.team); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetVCh() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_ShuffleTeam(t *testing.T) {
	type args struct {
		tchID string
		vss   []*dg.VoiceState
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.ShuffleTeam(tt.args.tchID, tt.args.vss); (err != nil) != tt.wantErr {
				t.Errorf("Manager.ShuffleTeam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_GetTeam(t *testing.T) {
	type args struct {
		tchID string
	}
	tests := []struct {
		name             string
		m                *Manager
		args             args
		wantTeam1UserIDs []string
		wantTeam2UserIDs []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTeam1UserIDs, gotTeam2UserIDs := tt.m.GetTeam(tt.args.tchID)
			if !reflect.DeepEqual(gotTeam1UserIDs, tt.wantTeam1UserIDs) {
				t.Errorf("Manager.GetTeam() gotTeam1UserIDs = %v, want %v", gotTeam1UserIDs, tt.wantTeam1UserIDs)
			}
			if !reflect.DeepEqual(gotTeam2UserIDs, tt.wantTeam2UserIDs) {
				t.Errorf("Manager.GetTeam() gotTeam2UserIDs = %v, want %v", gotTeam2UserIDs, tt.wantTeam2UserIDs)
			}
		})
	}
}

func TestManager_GetMatchByTChID(t *testing.T) {
	type args struct {
		tchID string
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		want    *Match
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetMatchByTChID(tt.args.tchID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetMatchByTChID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetMatchByTChID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetMatchStatus(t *testing.T) {
	type args struct {
		tchID string
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		want    *Status
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetMatchStatus(tt.args.tchID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetMatchStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetMatchStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_SetRecommendedChannel(t *testing.T) {
	type args struct {
		tchID string
		vch   *dg.Channel
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SetRecommendedChannel(tt.args.tchID, tt.args.vch); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetRecommendedChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_SetListeningMessage(t *testing.T) {
	type args struct {
		tchID string
		msg   *dg.Message
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SetListeningMessage(tt.args.tchID, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetListeningMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_IsListeningMessage(t *testing.T) {
	type args struct {
		tchID string
		msgID string
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsListeningMessage(tt.args.tchID, tt.args.msgID); got != tt.want {
				t.Errorf("Manager.IsListeningMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_getUsingTCh(t *testing.T) {
	type args struct {
		lock bool
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want []*dg.Channel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.getUsingTCh(tt.args.lock); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.getUsingTCh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_getUsingVCh(t *testing.T) {
	type args struct {
		lock bool
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want []*dg.Channel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.getUsingVCh(tt.args.lock); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.getUsingVCh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isContain(t *testing.T) {
	type args struct {
		s    string
		list []string
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
			if got := isContain(tt.args.s, tt.args.list); got != tt.want {
				t.Errorf("isContain() = %v, want %v", got, tt.want)
			}
		})
	}
}
