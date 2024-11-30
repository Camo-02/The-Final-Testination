package entity

import "backend/utils"

type Block struct {
	utils.Model

	Content string `gorm:"not null" json:"content"`

	// NULL when the block is not part of the answer, otherwise represents the index
	// of the block in the answer
	// It is of type `*uint` because `uint`'s null value is 0, which is a valid index, while `*uint`'s null value is nil
	Order *uint `json:"order"`

	// Whether the block is part of the skeleton
	Skeleton bool `json:"skeleton"`

	GameID string `gorm:"not null" json:"gameId"`
	Game   Game   `json:"-"`
}
