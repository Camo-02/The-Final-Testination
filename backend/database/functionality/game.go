package functionality

import (
	"backend/database"
	"backend/database/entity"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GameGetById(database *database.FinalTestinationDB, gameID uuid.UUID) (*entity.Game, error) {
	var game entity.Game
	var blocks []entity.Block
	result := database.Orm.Model(&entity.Game{}).Select("*").Where("id = ?", gameID).First(&game)
	res_blocks := database.Orm.Model(&entity.Block{}).Select("*").Where("game_id = ?", gameID).Scan(&blocks)

	err := result.Error
	err2 := res_blocks.Error
	if err != nil || err2 != nil {
		return nil, err
	}
	game.Blocks = blocks
	return &game, nil
}

func GameSetStartTime(database *database.FinalTestinationDB, gameID uuid.UUID, playerID uuid.UUID) (*entity.PlayerGame, error) {
	var playerGame *entity.PlayerGame

	result := database.Orm.First(&playerGame, "player_id = ? AND game_id = ?", playerID, gameID)

	if result.Error != nil {

		// if the player never played this game, then the playerGame will be nil. In this case, we create a new playerGame
		// with the provided values. Other values will be filled with default values by the DB.
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			playerGame.PlayerID = playerID.String()
			playerGame.GameID = gameID.String()
			playerGame.StartTime = time.Now()
			result = database.Orm.Create(&playerGame)
		}
	}

	return playerGame, result.Error
}
