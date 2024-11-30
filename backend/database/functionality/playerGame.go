package functionality

import (
	"backend/constants"
	"backend/database"
	"backend/database/entity"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaderboardEntry struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}

type GameHintsAvailability struct {
	FreezeTimeUsed   int `json:"freeze_time_spent_coins"`
	TextualHintUsed  int `json:"textual_hint_spent_coins"`
	SolutionHintUsed int `json:"solution_hint_spent_coins"`
}

func PlayerGetTotalCoins(db *database.FinalTestinationDB, playerID uuid.UUID) (int, error) {
	var totalCoins int

	// If the player hasn't completed any game, then the `Scan` call will not find any tuple, and the `totalCoins` will be 0 (default value)
	err := db.Orm.Model(&entity.PlayerGame{}).
		Select("SUM(score - textual_hint_points_used - hint_solution_points_used - time_freeze_points_used) as coins").
		Group("player_id").
		Having("player_id = ?", playerID).
		Scan(&totalCoins)

	return totalCoins, err.Error
}

func PlayerGameCreateMaxScore(db *database.FinalTestinationDB, gameID uuid.UUID, playerID uuid.UUID, end_time int64) (int, float64, error) {

	game_completion_data := struct {
		StartTime            int64  `json:"start_time"`
		PerfectTimeslot      int64  `json:"perfect_timeslot"`
		GreatTimeslot        int64  `json:"great_timeslot"`
		MediumTimeslot       int64  `json:"medium_timeslot"`
		NotSoGoodTimeslot    int64  `json:"not_so_good_timeslot"`
		MaxScore             int64  `json:"max_score"`
		WrongAttemptCost     int64  `json:"wrong_attempt_cost"`
		Attempts             int64  `json:"attempts"`
		TimeFreezeDuration   int64  `json:"time_freeze_duration"`
		TimeFreezePointsUsed *int64 `json:"time_freeze_points_used"`
	}{}

	tx := db.Orm.Model(&entity.PlayerGame{}).
		Joins("JOIN games ON player_games.game_id = games.id").
		Where("games.id = ? AND player_games.player_id = ?", gameID, playerID).
		Select("EXTRACT(EPOCH FROM player_games.start_time)::bigint as start_time, perfect_timeslot, great_timeslot, medium_timeslot, not_so_good_timeslot, max_score, wrong_attempt_cost, attempts, time_freeze_duration, time_freeze_points_used").
		First(&game_completion_data)
	log.Println("game_completion_data: ", game_completion_data)

	if tx.Error != nil {
		return 0, 0, tx.Error
	}

	time_used := end_time - game_completion_data.StartTime
	if game_completion_data.TimeFreezePointsUsed != nil {
		time_used -= game_completion_data.TimeFreezeDuration
	}
	var multiplier float64
	if time_used < game_completion_data.PerfectTimeslot {
		multiplier = 1
	} else if time_used < game_completion_data.GreatTimeslot {
		multiplier = 0.8
	} else if time_used < game_completion_data.MediumTimeslot {
		multiplier = 0.6
	} else if time_used < game_completion_data.NotSoGoodTimeslot {
		multiplier = 0.4
	} else {
		multiplier = 0.2
	}
	log.Println("time_used: ", time_used)
	log.Println("multiplier: ", multiplier)

	end_time_time := time.Unix(end_time, 0)

	score := int(math.Max(float64(game_completion_data.MaxScore)*multiplier-float64(game_completion_data.Attempts)*float64(game_completion_data.WrongAttemptCost), float64(game_completion_data.MaxScore)*0.2))

	log.Println("score: ", score)

	// update end time and score on the db
	tz := db.Orm.Model(&entity.PlayerGame{}).Where("game_id = ? AND player_id = ?", gameID, playerID).Updates(entity.PlayerGame{Score: score, EndTime: &end_time_time})

	return score, multiplier, tz.Error
}

func PlayerIncrementAttempts(db *database.FinalTestinationDB, gameID uuid.UUID, playerID uuid.UUID) error {
	tx := db.Orm.Exec("UPDATE player_games SET attempts = attempts + 1 WHERE player_id = ? AND game_id = ?", playerID, gameID)

	return tx.Error
}

func PlayerGameGet(database *database.FinalTestinationDB, playerID uuid.UUID, gameID uuid.UUID) (*entity.PlayerGame, error) {
	var playerGame entity.PlayerGame
	result := database.Orm.Where("user = ? AND game = ?", playerID, gameID).First(&playerGame)

	return &playerGame, result.Error
}

func GetLeaderboardPlayers(database *database.FinalTestinationDB, page int) (*[]LeaderboardEntry, error) {
	var playersScore []LeaderboardEntry

	//takes score as it is created in PlayerGameCreateMaxScore and subtracts hints used
	result := database.Orm.Model(&entity.PlayerGame{}).
		Joins("JOIN players ON player_games.player_id = players.id").
		Select("players.username, SUM(player_games.score - player_games.textual_hint_points_used - player_games.hint_solution_points_used - player_games.time_freeze_points_used) AS score").
		Group("players.id").
		Order("score DESC").
		Limit(constants.PAGE_SIZE).
		Offset((page * constants.PAGE_SIZE) - 1).
		Scan(&playersScore)

	return &playersScore, result.Error
}

func GetLeaderbordElementsNumber(database *database.FinalTestinationDB) (int, error) {
	var elemsNumber int

	result := database.Orm.Model(&entity.PlayerGame{}).Select("COUNT(DISTINCT player_id)").Scan(&elemsNumber)
	if result.Error != nil {
		return -1, result.Error
	}

	return elemsNumber, nil
}

func CheckGamePlayerCompleted(database *database.FinalTestinationDB, gameID uuid.UUID, playerID uuid.UUID) (bool, error) {
	var count int64
	result := database.Orm.Model(&entity.PlayerGame{}).Where("game_id = ? AND player_id = ? AND end_time IS NOT NULL", gameID, playerID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func PlayerGameUseHint(db *database.FinalTestinationDB, gameID uuid.UUID, playerID uuid.UUID, hintType string, order *int, playerCoins int) (int, *string, error) {
	var pg entity.PlayerGame
	err := db.Orm.First(&pg, "game_id = ? AND player_id = ?", gameID, playerID)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return 404, nil, err.Error
		}
		return 400, nil, err.Error
	}

	hintUsed := pg.TextualHintPointsUsed
	if hintType == "freeze" {
		hintUsed = pg.TimeFreezePointsUsed
	} else if hintType == "fill" {
		hintUsed = pg.HintSolutionPointsUsed
	}

	if hintUsed > 0 && pg.EndTime == nil {
		// Forbidden
		return 403, nil, fmt.Errorf("cannot buy hint with type %s a second time", hintType)
	}

	var game entity.Game
	err = db.Orm.Select("time_freeze_price, hint_solution_price, textual_hint_price, textual_hint, time_freeze_duration").First(&game, "id = ?", gameID)
	if err.Error != nil {
		return 500, nil, nil
	}

	var hintContent string
	if hintType == "textual" {
		if pg.EndTime == nil && playerCoins < game.TextualHintPrice {
			return 400, nil, fmt.Errorf("not enough coins")
		}
		hintContent = game.TextualHint
		pg.TextualHintPointsUsed = game.TextualHintPrice
	} else if hintType == "freeze" {
		if pg.EndTime == nil && playerCoins < game.TimeFreezePrice {
			return 400, nil, fmt.Errorf("not enough coins")
		}
		hintContent = fmt.Sprint(game.TimeFreezeDuration)
		pg.TimeFreezePointsUsed = game.TimeFreezePrice
	} else if hintType == "fill" {
		if pg.EndTime == nil && playerCoins < game.HintSolutionPrice {
			return 400, nil, fmt.Errorf("not enough coins")
		}
		res_blocks := db.Orm.Model(&entity.Block{}).Select("content").Where("game_id = ? AND \"order\" = ? AND skeleton = false", gameID, order).First(&hintContent)
		if res_blocks.Error != nil {
			if errors.Is(res_blocks.Error, gorm.ErrRecordNotFound) {
				return 400, nil, fmt.Errorf("block not found")

			}
			return 500, nil, fmt.Errorf("something went wrong, try again")

		}
		pg.HintSolutionPointsUsed = game.HintSolutionPrice
	}

	// Do not spend points if the game is already completed
	if pg.EndTime != nil {
		return 200, &hintContent, nil
	}

	//spend coins
	res := db.Orm.Where("game_id = ? AND player_id = ?", gameID, playerID).Updates(&pg)

	if res.Error != nil {
		return 500, nil, fmt.Errorf("couldn't update the hint usage")
	}

	return 200, &hintContent, nil
}
