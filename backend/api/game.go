package api

import (
	"backend/database"
	"backend/database/entity"
	"backend/database/functionality"
	"backend/middlewares"
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"errors"

	"gorm.io/gorm"
)

type hintUsedRequest struct {
	HintType string `json:"hint_type"`
	Order    *int   `json:"order"`
}

func SetUpGameRoutes(router *fiber.Router, database *database.FinalTestinationDB) {
	(*router).Get("/:gameId",
		middlewares.CheckValidUUID("gameId"),
		middlewares.InjectDB(database),
		middlewares.ValidateJWT,
		middlewares.CheckValidPlayer,
		getGame,
	)

	(*router).Post("/:gameId/hint",
		middlewares.CheckValidUUID("gameId"),
		middlewares.ParseBodyAsJSON[hintUsedRequest],
		middlewares.InjectDB(database),
		middlewares.ValidateJWT,
		middlewares.CheckValidPlayer,
		useHint,
	)
}

func getGame(c *fiber.Ctx) error {
	gameId := c.Locals("gameId").(uuid.UUID)
	db := c.Locals("db").(*database.FinalTestinationDB)
	player := c.Locals("player").(entity.Player)

	game, err := functionality.GameGetById(db, gameId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Couldn't find the game you're looking for",
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Couldn't get the game you're looking for",
			})
		}
	}

	// Check that the player has completed the previous game
	if game.GameOrder != 1 {
		previousGameOrder := game.GameOrder - 1

		var previousGameResult bool
		db.Orm.Raw("SELECT EXISTS( SELECT * FROM player_games WHERE player_id = ? AND game_id = (SELECT id FROM games WHERE game_order=?) AND end_time IS NOT NULL) AS found", player.ID, previousGameOrder).Scan(&previousGameResult)
		if !previousGameResult {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You have to complete the previous game first",
			})
		}
	}

	playerID := uuid.MustParse(player.ID)

	totalCoins, err := functionality.PlayerGetTotalCoins(db, playerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Couldn't get player's total coins",
		})
	}

	pg, err := functionality.GameSetStartTime(db, gameId, playerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Couldn't start the game",
		})
	}

	// if the user has already completed the level, signal to the frontend that the user can always
	// buy the hints, and that the hints are free
	if pg.EndTime != nil {
		pg.TimeFreezePointsUsed = 0
		pg.TextualHintPointsUsed = 0
		pg.HintSolutionPointsUsed = 0
		game.TimeFreezePrice = 0
		game.TextualHintPrice = 0
		game.HintSolutionPrice = 0
	}

	//TODO: should probably unit-test this part
	skeleton := fiber.Map{}
	blocks := []string{}
	solutionLength := 0

	for _, block := range game.Blocks {
		if block.Order != nil {
			solutionLength++
		}

		if block.Skeleton {
			// TODO: force block order to be != nil if block.Skeleton is true
			skeleton[fmt.Sprintf("%d", *block.Order)] = block.Content
		} else {
			blocks = append(blocks, block.Content)
		}
	}

	body := fiber.Map{
		"title":                game.Title,
		"story":                game.Story,
		"cheatsheet":           game.Cheatsheet,
		"skeleton":             skeleton,
		"blocks":               arrayShuffle(blocks),
		"solution_length":      solutionLength,
		"background":           game.Background,
		"winning_message":      game.WinningMessage,
		"player_coins":         totalCoins,
		"freeze_time_duration": game.TimeFreezeDuration,
		"freeze_time_price":    game.TimeFreezePrice,
		"freeze_time_used":     pg.TimeFreezePointsUsed,
		"textual_hint_price":   game.TextualHintPrice,
		"textual_hint_used":    pg.TextualHintPointsUsed,
		"solution_hint_price":  game.HintSolutionPrice,
		"solution_hint_used":   pg.HintSolutionPointsUsed,
	}

	return c.JSON(body)
}

// TODO: probably should be done in a better way
func arrayShuffle[T any](a []T) []T {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func useHint(c *fiber.Ctx) error {
	hintType := c.Locals("parsedBody").(hintUsedRequest).HintType
	if hintType != "freeze" && hintType != "textual" && hintType != "fill" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid hint type",
		})
	}

	var order *int
	if hintType == "fill" {
		order = c.Locals("parsedBody").(hintUsedRequest).Order
		if order == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "block not selected",
			})
		}
	}

	gameId := c.Locals("gameId").(uuid.UUID)
	db := c.Locals("db").(*database.FinalTestinationDB)
	player := c.Locals("player").(entity.Player)

	playerID := uuid.MustParse(player.ID)

	totalPoints, err := functionality.PlayerGetTotalCoins(db, playerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Coulnd't get the number of coins for the current user",
		})
	}

	statusCode, hintContent, err := functionality.PlayerGameUseHint(db, gameId, playerID, hintType, order, totalPoints)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(statusCode).JSON(fiber.Map{
				"error": "Couldn't find the game you're looking for",
			})
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"hintContent": hintContent})
}
