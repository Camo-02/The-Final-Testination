package main

import (
	"backend/api"
	"backend/database"
	"backend/env"
	"backend/loggers"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Ping(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func main() {
	app := fiber.New()

	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     "http://localhost:5173",
			AllowHeaders:     "Origin, Content-Type, Accept",
			AllowCredentials: true,
		},
	))

	app.Use(logger.New())

	app.Get("/ping", Ping)

	loggers.Info.Println("Creating database connection")
	db := database.CreateFinalTestinationDB(
		env.DB_USERNAME, env.DB_PASSWORD,
		env.DB_HOST, env.DB_PORT, env.DB_NAME)
	db.CreateSchemas()

	blockRouter := app.Group("/blocks")
	api.SetUpBlocksRoutes(&blockRouter, db)

	gameRouter := app.Group("/game")
	api.SetUpGameRoutes(&gameRouter, db)

	leaderboardRouter := app.Group("/leaderboard")
	api.SetUpPlayerGameRoutes(&leaderboardRouter, db)

	playerRouter := app.Group("/player")
	api.SetUpPlayerRoutes(&playerRouter, db)

	port := fmt.Sprintf(":%s", env.API_PORT)
	loggers.Error.Fatal(app.Listen(port))

}
