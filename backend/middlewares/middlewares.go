package middlewares

import (
	"backend/constants"
	"backend/database"
	"backend/database/functionality"
	"backend/jwt"
	"backend/loggers"
	"backend/validators"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func InjectDB(db *database.FinalTestinationDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	}
}

func CheckValidUUID(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid, err := uuid.Parse(c.Params(key))
		if err != nil {
			//TODO: verify error log
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Please provide a valid UUID representing a game",
			})
		}
		c.Locals(key, uuid)
		return c.Next()
	}
}

func ParseBodyAsJSON[T any](c *fiber.Ctx) error {
	var body T

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Please provide a body that matches the expected structure",
		})
	}

	c.Locals("parsedBody", body)

	return c.Next()
}

func CheckValidPageNumber(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.Params(key)
		int_page, err := strconv.Atoi(page)
		if err != nil {
			int_page = 1
		}

		if int_page < 1 {
			int_page = 1
		}

		c.Locals("page", int_page)
		return c.Next()
	}
}

func ValidateBodyAs[T validators.Valid](c *fiber.Ctx) error {
	body := c.Locals("parsedBody")

	if err := body.(T).Validate(validators.Validate.GetValidator()); err != nil {
		loggers.Debug.Printf("Error validating body (%+v): %v", body, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "wrong credentials"})
	}
	return c.Next()
}

func ValidateJWT(c *fiber.Ctx) error {
	cookie := c.Cookies(constants.AUTH_COOKIE_NAME)

	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing cookie"})
	}

	testinationClaims, err := jwt.ParseJWT(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	c.Locals("claims", testinationClaims)
	return c.Next()
}

func CheckValidPlayer(c *fiber.Ctx) error {
	db := c.Locals("db").(*database.FinalTestinationDB)
	claims := c.Locals("claims").(*jwt.FinalTestinationClaims)

	player, err := functionality.PlayerGetByID(db, claims.PlayerID)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	c.Locals("player", *player)
	return c.Next()
}
