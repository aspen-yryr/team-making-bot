package migration

import (
	tm "github.com/aspen-yryr/team-making-bot/pkg/team-maker"

	"gorm.io/gorm"
)

type Migrator interface {
	Run() error
}

type SampleMigrator struct {
	db *gorm.DB
}

func (s SampleMigrator) Run() error {
	s.db.AutoMigrate(&tm.BattleRecord{})
	return nil
}

func NewSampleMigrator(db *gorm.DB) Migrator {
	return &SampleMigrator{db}
}
