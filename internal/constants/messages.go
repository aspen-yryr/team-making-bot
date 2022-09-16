package constants

import (
	"fmt"
)

var Msgs *Messages

// func init() {
// 	Msgs = NewMessages()
// }

type Message struct {
	format *string
	MessageFormatter
}

type MessageFormatter interface {
	Format([]string) string
}

func (f *Message) Format(args ...interface{}) string {
	// TODO: required arguments check
	if len(args) == 0 {
		return fmt.Sprintf(*f.format)
	}
	return fmt.Sprintf(*f.format, args...)
}

func newMessage(format string) *Message {
	return &Message{
		format: &format,
	}
}

const help = `- 操作方法
このBotへのメンション、リプ、スタンプで操作
チーム分け機能はテキストチャンネルごとに独立しています。
- コマンド
start：チーム分けの開始
end：チーム分けの終了
reset：すべての設定をリセット
help：このヘルプの表示`
const exit = `このチャンネルでのチーム分けを終了します。`
const unknownError = `エラーが発生しました。`
const matchAlreadyStarted = `このテキストチャンネルではすでにチーム分けが開始しています。`
const askTeam1VCh = `チームAが使用するボイスチャンネルは？`
const askTeam1VChWithRecommend = `チームAが使用するボイスチャンネルは%sでいいですか？`
const askTeam2VCh = `チームBが使用するボイスチャンネルは？`
const confirmTeam1VCh = `チームAは%sチャンネルを使用します。`
const confirmTeam2VCh = `チームBは%sチャンネルを使用します。`
const requestChName = `チャンネル名を入力してください。`
const noVChAvailable = `使用できるボイスチャンネルがありません。`
const makeTeam = `ボイスチャンネル%sと%sにいるメンバーをチーム分けします。`
const confirmTeam = `このチーム分けで良いですか？
Team A(%s): %v
Team B(%s): %v
go：ボイスチャンネルを移動
shuffle：チーム分けやり直し`
const conflictVCh = `%sは使用できません。`
const ownerNotInVchs = `startコマンドを入力した人がvcに参加している必要があります。
vcに参加してからやり直してください。`

type Messages struct {
	Help,
	Exit,
	UnknownError,
	MatchAlreadyStarted,
	AskTeam1VCh,
	AskTeam1VChWithRecommend,
	AskTeam2VCh,
	ConfirmTeam1VCh,
	ConfirmTeam2VCh,
	RequestChName,
	MakeTeam,
	ConfirmTeam,
	ConflictVCh,
	NoVChAvailable,
	OwnerNotInVchs *Message
}

func NewMessages() *Messages {
	ms := &Messages{
		Help:                     newMessage(help),
		Exit:                     newMessage(exit),
		UnknownError:             newMessage(unknownError),
		MatchAlreadyStarted:      newMessage(matchAlreadyStarted),
		AskTeam1VCh:              newMessage(askTeam1VCh),
		AskTeam1VChWithRecommend: newMessage(askTeam1VChWithRecommend),
		AskTeam2VCh:              newMessage(askTeam2VCh),
		ConfirmTeam1VCh:          newMessage(confirmTeam1VCh),
		ConfirmTeam2VCh:          newMessage(confirmTeam2VCh),
		RequestChName:            newMessage(requestChName),
		MakeTeam:                 newMessage(makeTeam),
		ConfirmTeam:              newMessage(confirmTeam),
		ConflictVCh:              newMessage(conflictVCh),
		NoVChAvailable:           newMessage(noVChAvailable),
		OwnerNotInVchs:           newMessage(ownerNotInVchs),
	}
	return ms
}
