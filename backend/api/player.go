package api

import (
	"backend/constants"
	"backend/database"
	"backend/database/entity"
	"backend/database/functionality"
	"backend/jwt"
	"backend/middlewares"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type userCredentials struct {
	Credential string `binding:"required" validate:"testination-credential"`
	Password   string `binding:"required" validate:"testination-password"`
}

type googleLoginRequest struct {
	Token string `json:"token" binding:"required"`
}

func (u userCredentials) Validate(v *validator.Validate) error {
	return v.Struct(u)
}

type userRegister struct {
	Username string `binding:"required" validate:"testination-username"`
	Password string `binding:"required" validate:"testination-password"`
	Email    string `binding:"required" validate:"email"`
}

type icon struct {
	Icon string `binding:"required" json:"icon"`
}

func (u userRegister) Validate(v *validator.Validate) error {
	return v.Struct(u)
}

func SetUpPlayerRoutes(router *fiber.Router, database *database.FinalTestinationDB) {
	(*router).Post("/register", middlewares.ParseBodyAsJSON[userRegister], middlewares.ValidateBodyAs[userRegister], middlewares.InjectDB(database), register)
	(*router).Post("/login", middlewares.ParseBodyAsJSON[userCredentials], middlewares.ValidateBodyAs[userCredentials], middlewares.InjectDB(database), login)
	(*router).Post("/loggedInfo", middlewares.InjectDB(database), middlewares.ValidateJWT, middlewares.CheckValidPlayer, getLoggedInfo)
	(*router).Get("/availableLevels", middlewares.InjectDB(database), middlewares.ValidateJWT, middlewares.CheckValidPlayer, seeAvailableLevels)
	(*router).Get("/profile", middlewares.InjectDB(database), middlewares.ValidateJWT, middlewares.CheckValidPlayer, getProfile)
	(*router).Get("/availableIcons", middlewares.InjectDB(database), middlewares.ValidateJWT, middlewares.CheckValidPlayer, getAvailableIcons)
	(*router).Post("/changeIcon", middlewares.ParseBodyAsJSON[icon], middlewares.InjectDB(database), middlewares.ValidateJWT, middlewares.CheckValidPlayer, changeIcon)
	(*router).Post("/logout", logOut)
	(*router).Post("/googleLogin", middlewares.ParseBodyAsJSON[googleLoginRequest], middlewares.InjectDB(database), loginWithGoogle) // new endpoint for Google login
}

func register(c *fiber.Ctx) error {
	db := c.Locals("db").(*database.FinalTestinationDB)
	body := c.Locals("parsedBody").(userRegister)

	user, err := functionality.PlayerCreate(db, body.Username, body.Password, body.Email)
	if err != nil {
		// TODO: diversify errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Impossible to create a new user": err.Error()})
	}

	return c.JSON(user)
}

func seeAvailableLevels(c *fiber.Ctx) error {
	db := c.Locals("db").(*database.FinalTestinationDB)
	player := c.Locals("player").(entity.Player)

	result, err := functionality.GetPlayerLevels(db, player.ID)
	if err != nil {
		// TODO: diversify errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func login(c *fiber.Ctx) error {
	db := c.Locals("db").(*database.FinalTestinationDB)
	userAnswer := c.Locals("parsedBody").(userCredentials)

	var user *entity.Player
	var err error

	if strings.Contains(userAnswer.Credential, "@") {
		user, err = functionality.PlayerGetByEmailAndPassword(db, userAnswer.Credential, userAnswer.Password)
	} else {
		user, err = functionality.PlayerGetByUsernameAndPassword(db, userAnswer.Credential, userAnswer.Password)
	}

	if err != nil {
		// TODO: could also be another error
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "wrong credentials",
		})
	}

	token, err := jwt.GenerateJWT(user.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not generate token",
		})
	}

	const JWT_EXPIRATION int = 60 * 60 * 24 * 31

	c.Cookie(&fiber.Cookie{
		Name:     constants.AUTH_COOKIE_NAME,
		Value:    token,
		MaxAge:   JWT_EXPIRATION,
		SameSite: "None",
		Secure:   true,
		HTTPOnly: true,
	})

	return c.JSON(user)
}

// Function for Google login
func loginWithGoogle(c *fiber.Ctx) error {
	fmt.Println("Request received at /googleLogin")
	db := c.Locals("db").(*database.FinalTestinationDB)
	body := c.Locals("parsedBody").(googleLoginRequest) // body.Token contains the Google token

	// 1. Google token verification
	googleAPIURL := "https://oauth2.googleapis.com/tokeninfo?id_token=" + body.Token
	resp, err := http.Get(googleAPIURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Google token",
		})
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var googleData struct {
		Iss           string `json:"iss"`            // Issuer
		Aud           string `json:"aud"`            // Audience
		Email         string `json:"email"`          // Email
		EmailVerified string `json:"email_verified"` // "true" or "false" as string
		Name          string `json:"name"`           // Full name
		Picture       string `json:"picture"`        // Profile picture URL
		Exp           string `json:"exp"`            // Expiration time
		Iat           string `json:"iat"`            // Issued at time
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse Google API response",
		})
	}

	// Verify the token is for your Client ID
	if googleData.Aud != "Your Client ID" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token audience mismatch",
		})
	}

	// Check if the email is verified
	if googleData.EmailVerified == "false" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Google email not verified",
		})
	}

	// Use the data to authenticate or register the user
	fmt.Printf("Email: %s, Name: %s\n", googleData.Email, googleData.Name)

	// 2. Search for the user in the database
	user, err := functionality.PlayerGetByEmail(db, googleData.Email)
	if err != nil {
		// Handle database errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve user",
		})
	}

	if user == nil {
		// User not found, create it
		user, err = functionality.PlayerCreate(db, googleData.Name, "", googleData.Email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "could not create user",
			})
		}
	} else {
		// User found, update data if necessary
		user.Username = googleData.Name
		if err := functionality.PlayerUpdate(db, user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "could not update user",
			})
		}
	}

	// 3. Generate a JWT token
	token, err := jwt.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not generate token",
		})
	}

	const JWT_EXPIRATION int = 60 * 60 * 24 * 31
	c.Cookie(&fiber.Cookie{
		Name:     constants.AUTH_COOKIE_NAME,
		Value:    token,
		MaxAge:   JWT_EXPIRATION,
		SameSite: "None",
		Secure:   true,
		HTTPOnly: true,
	})

	return c.JSON(user)
}

func getLoggedInfo(c *fiber.Ctx) error {
	player := c.Locals("parsedBody").(*entity.Player)

	return c.JSON(player)
}

func getProfile(c *fiber.Ctx) error {
	db := c.Locals("db").(*database.FinalTestinationDB)
	player := c.Locals("player").(entity.Player)
	result, err := functionality.Profile(db, &player)
	if err != nil {
		// TODO: diversify errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func getAvailableIcons(c *fiber.Ctx) error {
	db := c.Locals("db").(*database.FinalTestinationDB)

	icons, err := functionality.GetAvailableIcons(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(icons)
}

func changeIcon(c *fiber.Ctx) error {
	db := c.Locals("db").(*database.FinalTestinationDB)
	player := c.Locals("player").(entity.Player)
	body := c.Locals("parsedBody").(icon)
	err := functionality.ChangeIcon(db, player.ID, body.Icon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusOK)
}
func logOut(c *fiber.Ctx) error {
	fmt.Println("Logout endpoint called")

	// Clear the cookie
	c.ClearCookie(constants.AUTH_COOKIE_NAME)

	// Explicitly set a cookie with the same properties but with expiration in the past
	c.Cookie(&fiber.Cookie{
		Name:     constants.AUTH_COOKIE_NAME,
		Value:    "",
		MaxAge:   -1, // Set to a negative value
		SameSite: "None",
		Secure:   true, // Keep true in production
		HTTPOnly: true,
		Path:     "/", // Match the cookie path to ensure deletion
	})

	fmt.Println("Cookie cleared")
	return nil
}
