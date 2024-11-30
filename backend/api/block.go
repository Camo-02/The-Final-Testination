package api

import (
	"backend/database"
	"backend/database/entity"
	"backend/database/functionality"
	"backend/middlewares"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"errors"

	"gorm.io/gorm"
)

type blockAnswer struct {
	Blocks []string `json:"blocks"`
}

func SetUpBlocksRoutes(router *fiber.Router, database *database.FinalTestinationDB) {
	(*router).Post("/:gameId/check-answer",
		middlewares.CheckValidUUID("gameId"),
		middlewares.ParseBodyAsJSON[blockAnswer],
		middlewares.ValidateJWT,
		middlewares.InjectDB(database),
		middlewares.CheckValidPlayer,
		checkAnswer,
	)
}

func checkAnswer(c *fiber.Ctx) error {
	gameId := c.Locals("gameId").(uuid.UUID)
	db := c.Locals("db").(*database.FinalTestinationDB)
	blockAnswer := c.Locals("parsedBody").(blockAnswer)
	player := c.Locals("player").(entity.Player)
	playerID, err := uuid.Parse(player.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid player ID",
		})
	}

	solution, err := functionality.BlocksOfAnswer(db, gameId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Couldn't find the game you're looking for",
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Couldn't check your answer, please try again later",
			})
		}
	}

	correct_indexes, is_all_correct := ArraysMatch(solution, blockAnswer.Blocks)
	completed, err := functionality.CheckGamePlayerCompleted(db, gameId, playerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Couldn't check if the game is completed",
		})
	}
	if !is_all_correct {
		if !completed {
			err := functionality.PlayerIncrementAttempts(db, gameId, playerID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed increment attempts",
				})
			}
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"matches": correct_indexes,
		})
	}

	var score int
	var multiplier float64
	if !completed {
		// TODO: should return also time_slot
		score, multiplier, err = functionality.PlayerGameCreateMaxScore(db, gameId, playerID, time.Now().Unix())
		if err != nil {
			//TODO: check for specific errors
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot update the score for this game.",
			})

		}

	}

	next_gameID, err := functionality.GameGetByPreviousGame(db, gameId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Couldn't get the next game",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"answer_correctly": "correct answer",
		"next_level_id":    next_gameID,
		"score":            score,
		"multiplier":       multiplier,
	})
}

func ArraysMatch(solution []string, answer []string) ([]bool, bool) {
	correct_indexes := make([]bool, len(solution))

	if len(solution) != len(answer) {
		return correct_indexes, false
	}

	is_all_correct := true

	for i, v := range answer {
		correct_indexes[i] = solution[i] == v
		if !correct_indexes[i] {
			is_all_correct = false
		}
	}

	return correct_indexes, is_all_correct
}
