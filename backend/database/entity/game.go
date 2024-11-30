package entity

import (
	"backend/utils"
)

type Game struct {
	utils.Model
	Title              string       `gorm:"not null" json:"title"`
	GameOrder          int          `gorm:"not null" json:"game_order"`
	Blocks             []Block      `json:"blocks,omitempty"`
	Story              string       `gorm:"not null" json:"story"`
	Cheatsheet         string       `gorm:"not null" json:"cheatsheet"`
	MaxScore           int          `gorm:"not null" json:"max_score" default:"0"`
	Description        string       `gorm:"not null" json:"description"`
	Background         string       `gorm:"not null" json:"background"`
	WinningMessage     string       `gorm:"not null" json:"winning_message"`
	PlayerGames        []PlayerGame `json:"playerGames,omitempty"`
	WrongAttemptCost   int          `gorm:"not null" json:"wrong_attempt_cost"`
	PerfectTimeslot    int          `gorm:"not null" json:"perfect_timeslot"`     // In seconds
	GreatTimeslot      int          `gorm:"not null" json:"great_timeslot"`       // In seconds
	MediumTimeslot     int          `gorm:"not null" json:"medium_timeslot"`      // In seconds
	NotSoGoodTimeslot  int          `gorm:"not null" json:"not_so_good_timeslot"` // In seconds
	TextualHintPrice   int          `gorm:"not null" json:"textual_hint_price"`
	TextualHint        string       `gorm:"not null" json:"textual_hint"`
	HintSolutionPrice  int          `gorm:"not null" json:"hint_solution_price"`
	TimeFreezePrice    int          `gorm:"not null" json:"time_freeze_price"`
	TimeFreezeDuration int          `gorm:"not null" json:"time_freeze_duration"`
}
