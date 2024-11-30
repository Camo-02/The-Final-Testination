package entity

import "backend/utils"

type Icon struct {
	utils.Model
	Svg string `gorm:"not null" json:"svg"`
}
