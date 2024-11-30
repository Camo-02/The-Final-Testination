package api

import (
	"backend/constants"
	"backend/database"
	"backend/database/functionality"
	"backend/middlewares"
	"math"

	"github.com/gofiber/fiber/v2"
)

func SetUpPlayerGameRoutes(router *fiber.Router, database *database.FinalTestinationDB) {
	(*router).Get("/:page", middlewares.CheckValidPageNumber("page"), middlewares.InjectDB(database), getPlayersData)
}

func getPlayersData(c *fiber.Ctx) error {
	page := c.Locals("page").(int)
	db := c.Locals("db").(*database.FinalTestinationDB)

	elementCount, err := functionality.GetLeaderbordElementsNumber(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"queryError": "An error occurred while fetching the number of pages" + err.Error(),
		})
	}

	maxPageNumber := int(math.Ceil(float64(elementCount) / constants.PAGE_SIZE))

	if page > maxPageNumber {
		page = maxPageNumber
	}

	res, err := functionality.GetLeaderboardPlayers(db, (page - 1))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"queryError": "An error occurred while fetching the leaderboard data" + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"currentPage": page,
		"pages":       maxPageNumber,
		"entries":     res,
	})
}
