package team_maker

import (
	"time"

	"gorm.io/gorm"
)

type BattleRecord struct {
	gorm.Model
	BattleResult
	BattleDateTime time.Time
}

type BattleResult struct {
	PlayerID string
	Kill     int
	Death    int
	Assist   int
	Victory  bool
}

type BattleRecordRepository interface {
	GetById(string) (*BattleRecord, error)
	FindByPlayerId(string) ([]BattleRecord, error)
	Save(*BattleRecord) error
}

type BattleRecordTeamMaker struct {
	*BaseTeamMaker
	repo BattleRecordRepository
}

type battleRecordScorer struct {
	repo BattleRecordRepository
}

func (sc battleRecordScorer) GetScore(id string) (float64, error) {
	ps, err := sc.repo.FindByPlayerId(id)
	if err != nil {
		return 0, err
	}

	if len(ps) == 0 {
		return 0, nil
	}

	won := 0
	for _, p := range ps {
		if p.Victory {
			won += 1
		}
	}

	return float64(won / len(ps)), nil
}

func (tm *BattleRecordTeamMaker) SetPlayers(players []string) {
	tm.DiscordIds = players
}

func (tm *BattleRecordTeamMaker) RegisterBattleResult(result *BattleResult) error {
	return tm.repo.Save(&BattleRecord{BattleResult: *result, BattleDateTime: time.Now()})
}

func NewBattleRecordTeamMaker(repo BattleRecordRepository) *BattleRecordTeamMaker {
	return &BattleRecordTeamMaker{
		NewTeamMaker(nil, battleRecordScorer{repo}),
		repo,
	}
}
