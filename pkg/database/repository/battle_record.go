package repository

import (
	"fmt"

	tm "github.com/aspen-yryr/team-making-bot/pkg/team-maker"

	"gorm.io/gorm"
)

type BattleRecordRepository struct {
	db *gorm.DB
}

func NewBattleRecordRepository(db *gorm.DB) *BattleRecordRepository {
	return &BattleRecordRepository{db}
}

func (r *BattleRecordRepository) GetById(id string) (*tm.BattleRecord, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *BattleRecordRepository) FindByPlayerId(id string) ([]tm.BattleRecord, error) {
	var battle_records []tm.BattleRecord
	r.db.Where("player_id = ?", id).Find(&battle_records)
	return battle_records, nil
}

func (r *BattleRecordRepository) Save(record *tm.BattleRecord) error {
	result := r.db.Create(record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
