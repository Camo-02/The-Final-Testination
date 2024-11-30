package functionality

import (
	"backend/database"
	"backend/database/entity"
	"backend/utils"
	"errors"

	"github.com/google/uuid"
)

func BlocksByGameID(database *database.FinalTestinationDB, gameID uuid.UUID) ([]entity.Block, error) {
	var blocks []entity.Block
	result := database.Orm.Where("game_id = ?", gameID).Find(&blocks)
	if result.Error != nil {
		return nil, errors.New("database error")
	}
	return blocks, nil
}

func BlocksOfAnswer(database *database.FinalTestinationDB, gameID uuid.UUID) ([]string, error) {
	var blocks []entity.Block
	result := database.Orm.Where("game_id = ? AND \"order\" IS NOT NULL", gameID).Order("\"order\"").Select("content").Find(&blocks)
	if result.Error != nil {
		return nil, result.Error
	}

	return utils.Map(blocks, func(b entity.Block) string { return b.Content }), nil
}

func GameGetByPreviousGame(database *database.FinalTestinationDB, previousGameId uuid.UUID) (uuid.UUID, error) {
	var id_string string
	result := database.Orm.Model(&entity.Game{}).Select("id").Where("game_order = (SELECT game_order FROM games WHERE id = ?)+1", previousGameId).First(&id_string)
	id, _ := uuid.Parse(id_string)
	return id, result.Error

	/*
		var nextGame entity.Game
		result := database.Orm.Model(&entity.Game{}).Where("game_order = (SELECT game_order FROM games WHERE id = ?)+1", previousGameId).First(&nextGame)
		log.Println("RESULT:", nextGame, "\tERROR:", result.Error)
		id,_ := uuid.Parse(nextGame.ID)
		return id, result.Error*/
}
