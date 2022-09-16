package repository

import (
	"fmt"

	"github.com/aspen-yryr/team-making-bot/internal/user"

	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return UserRepository{db}
}

type UserRepository struct {
	db *gorm.DB
}

func (mr UserRepository) Get(id string) (*user.User, error) {
	return nil, fmt.Errorf("not implemented")
}
