package entity

import "time"

// problem: hints that are free are not distinguishable from unused ones with this implementation

// Solution -> points used become *int (remove not null and default) and then:
//
//	if the points used is nil, then the hint was not used
//	if the points used is not nil, then the hint was used and was free
//	if the points used has another value, then the hint was used and was not free

type PlayerGame struct {
	PlayerID               string     `gorm:"not null; uniqueIndex:idx_gameid_playerid" json:"player_id"`
	Player                 Player     `json:"-"`
	GameID                 string     `gorm:"not null; uniqueIndex:idx_gameid_playerid" json:"gameId"`
	Game                   Game       `json:"-"`
	Score                  int        `gorm:"" json:"score"`
	Attempts               int        `gorm:"not null; default:0" json:"attempts"`
	StartTime              time.Time  `gorm:"not null; default:CURRENT_TIMESTAMP" json:"start_time"`
	EndTime                *time.Time `gorm:"" json:"end_time"`
	TextualHintPointsUsed  int        `gorm:"not null; default:0" json:"textual_hint_points_used"`
	HintSolutionPointsUsed int        `gorm:"not null; default:0" json:"hint_solution_points_used"`
	TimeFreezePointsUsed   int        `gorm:"not null; default:0" json:"time_freeze_points_used"`
}
