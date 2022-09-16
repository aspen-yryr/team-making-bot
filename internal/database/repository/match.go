package repository

import (
	"fmt"

	"github.com/aspen-yryr/team-making-bot/service/match"

	"gorm.io/gorm"
)

type Status int

const (
	StateVCh1Setting Status = iota
	StateVCh2Setting
	StateTeamPreview
)

type Match struct {
	gorm.Model
	OwnerID string
	Team1   []string
	Team2   []string
	Status  Status
}

type MatchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) match.MatchRepository {
	return MatchRepository{db}
}

func (mr MatchRepository) Create() (*match.Match, error) {
	return nil, fmt.Errorf("not implemented")
}

func (mr MatchRepository) Get(mt *match.Match) (*match.Match, error) {
	return nil, fmt.Errorf("not implemented")
}
