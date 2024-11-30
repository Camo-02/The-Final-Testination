package entity

import "backend/utils"

type Player struct {
	utils.Model
	Username    string       `gorm:"not null;unique" json:"username"`
	Password    string       `gorm:"not null" json:"-"`
	Email       string       `gorm:"not null;unique" json:"email"`
	PlayerGames []PlayerGame `json:"playerGames,omitempty"`
	IconID      string       `json:"iconId"`
	Secure      bool         `json:"secure"`
	SameSite    string       `json:"sameSite"`
}
