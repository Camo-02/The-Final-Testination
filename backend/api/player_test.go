package api

import (
	"backend/constants"
	"backend/database"
	"backend/database/functionality"
	"backend/env"
	"backend/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetLogin(t *testing.T) {
	tests := []struct {
		description  string // description of the test case
		credential   string
		password     string
		expectedCode int
	}{
		{
			description:  "Log in with a user that exists (username & password)",
			expectedCode: 200,
			credential:   "admin",
			password:     "rootroot",
		},
		{
			description:  "Log in with a user that exists (email & password)",
			expectedCode: 200,
			credential:   "admin@testination.com",
			password:     "rootroot",
		},
		{
			description:  "Log in with a user that doesn't exist",
			expectedCode: 401,
			credential:   "NOTadmin",
			password:     "rootroot",
		},
		{
			description:  "Log in with a wrong password",
			expectedCode: 401,
			credential:   "admin",
			password:     "wrongpassword",
		},
	}

	db := database.CreateFinalTestinationDB(
		env.DB_USERNAME, env.DB_PASSWORD,
		env.DB_HOST, env.DB_PORT, env.DB_NAME)
	db.CreateSchemas()

	app := fiber.New()
	playerGroup := app.Group("/player")
	SetUpPlayerRoutes(&playerGroup, db)

	for _, test := range tests {
		resp := utils.MockLogin(t, app, test.credential, test.password)
		assert.Equal(t, test.expectedCode, resp.StatusCode, test.description)
		cookie := utils.Filter(resp.Cookies(), func(c *http.Cookie) bool {
			return c.Name == constants.AUTH_COOKIE_NAME
		})
		if test.expectedCode == 200 {
			assert.NotEmpty(t, cookie, test.description)
			assert.NotEmpty(t, cookie[0].Value, test.description)
		} else {
			assert.Empty(t, cookie, test.description)
		}
	}
}

var availableLevelsRoute string = "/player/availableLevels"

var allGamesCompletedExpectedBodyAsString string = `[
    {
        "game_id": "af8e4754-1b84-4fec-bec4-154a3f894b8f",
		"title": "XSS Attack",
        "game_order": 1,
        "max_score": 100,
		"score_achieved": 100,
        "description": "A conspiracy hidden behind a popular corporation aims to use pollution to profit off the misfortune of countless people. Time to find more about it."
    },
    {
        "game_id": "05732286-9fa5-45d4-bef3-13ae0d481afa",
		"title": "SQL Injection",
        "game_order": 2,
        "max_score": 200,
		"score_achieved": 200,
        "description": "Your target is traveling to meet the other conspirators: find where he's going and you'll get to the source of the problem."
    },
    {
        "game_id": "a76db50b-ee98-4dd4-9d63-4c0ab695ad5f",
		"title": "",
        "game_order": 3,
        "max_score": 300,
		"score_achieved": 300,
        "description": ""
    }
]`

var noGamesCompletedExpectedBodyAsString string = `[
    {
        "game_id": "af8e4754-1b84-4fec-bec4-154a3f894b8f",
		"title": "XSS Attack",
        "game_order": 1,
        "max_score": 100,
		"score_achieved": 0,
        "description": "A conspiracy hidden behind a popular corporation aims to use pollution to profit off the misfortune of countless people. Time to find more about it."
    },
    {
        "game_id": "05732286-9fa5-45d4-bef3-13ae0d481afa",
		"title": "SQL Injection",
        "game_order": 2,
        "max_score": -1,
		"score_achieved": 0,
        "description": "locked"
    },
    {
        "game_id": "a76db50b-ee98-4dd4-9d63-4c0ab695ad5f",
		"title": "",
        "game_order": 3,
        "max_score": -1,
		"score_achieved": 0,
        "description": "locked"
    }
]`

var someGamesCompletedExpectedBodyAsString string = `[
	{
	  "game_id": "af8e4754-1b84-4fec-bec4-154a3f894b8f",
	  "title": "XSS Attack",
	  "game_order": 1,
	  "max_score": 100,
	  "score_achieved": 90,
	  "description": "A conspiracy hidden behind a popular corporation aims to use pollution to profit off the misfortune of countless people. Time to find more about it."
	},
	{
	  "game_id": "05732286-9fa5-45d4-bef3-13ae0d481afa",
	  "title": "SQL Injection",
	  "game_order": 2,
	  "max_score": 200,
	  "score_achieved": 0,
	  "description": "Your target is traveling to meet the other conspirators: find where he's going and you'll get to the source of the problem."
	},
	{
	  "game_id": "a76db50b-ee98-4dd4-9d63-4c0ab695ad5f",
	  "title": "",
	  "game_order": 3,
	  "max_score": -1,
	  "score_achieved": 0,
	  "description": "locked"
	}
  ]`

func TestGetAvailableLevels(t *testing.T) {
	tests := []struct {
		method       string
		description  string // description of the test case
		username     string
		password     string
		route        string // route path to test
		expectedCode int    // expected HTTP status code
		body         string
		expected     string
	}{
		{
			method:       "GET",
			description:  "Player completed all levels",
			username:     "admin",
			password:     "rootroot",
			route:        availableLevelsRoute,
			expectedCode: 200,
			expected:     allGamesCompletedExpectedBodyAsString,
		},
		{
			method:       "GET",
			description:  "Player completed some level",
			username:     "test",
			password:     "rootroot",
			route:        availableLevelsRoute,
			expectedCode: 200,
			expected:     someGamesCompletedExpectedBodyAsString,
		},
		{
			method:       "GET",
			description:  "Player completed no level",
			username:     "newUser",
			password:     "rootroot",
			route:        availableLevelsRoute,
			expectedCode: 200,
			expected:     noGamesCompletedExpectedBodyAsString,
		},
	}

	db := database.CreateFinalTestinationDB(
		env.DB_USERNAME, env.DB_PASSWORD,
		env.DB_HOST, env.DB_PORT, env.DB_NAME)
	db.CreateSchemas()

	app := fiber.New()
	playerGroup := app.Group("/player")
	SetUpPlayerRoutes(&playerGroup, db)

	for _, test := range tests {
		loginResp := utils.MockLogin(t, app, test.username, test.password)
		cookie := utils.Filter(loginResp.Cookies(), func(c *http.Cookie) bool {
			return c.Name == constants.AUTH_COOKIE_NAME
		})[0].Value
		req := httptest.NewRequest(test.method, test.route, bytes.NewBuffer([]byte(test.body)))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: constants.AUTH_COOKIE_NAME, Value: cookie})

		resp, err := app.Test(req, -1) // -1 means no timeout

		assert.NoError(t, err)
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		responseBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		//TODO: should probably check the error message
		if resp.StatusCode != 200 {
			continue
		}

		fmt.Println(string(responseBody))
	}
}

var profileRoute string = "/player/profile"

func TestGetProfile(t *testing.T) {
	tests := []struct {
		description string
		username    string
		password    string
	}{
		{
			description: "Get profile of a user that exists",
			username:    "admin",
			password:    "rootroot",
		},
	}

	db := database.CreateFinalTestinationDB(
		env.DB_USERNAME, env.DB_PASSWORD,
		env.DB_HOST, env.DB_PORT, env.DB_NAME)
	db.CreateSchemas()

	app := fiber.New()
	gameGroup := app.Group("/game")
	SetUpGameRoutes(&gameGroup, db)
	playerGroup := app.Group("/player")
	SetUpPlayerRoutes(&playerGroup, db)

	for _, test := range tests {
		loginResp := utils.MockLogin(t, app, test.username, test.password)
		cookie := utils.Filter(loginResp.Cookies(), func(c *http.Cookie) bool {
			return c.Name == constants.AUTH_COOKIE_NAME
		})[0].Value
		req := httptest.NewRequest("GET", profileRoute, nil)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: constants.AUTH_COOKIE_NAME, Value: cookie})

		resp, errTest := app.Test(req, -1)
		assert.NoError(t, errTest)
		var profile functionality.ProfileDTO
		body, errBody := io.ReadAll(resp.Body)
		assert.NoError(t, errBody)
		errUnmarshal := json.Unmarshal(body, &profile)
		assert.NoError(t, errUnmarshal)
		assert.Equal(t, 200, resp.StatusCode, test.description)
		assert.Equal(t, test.username, profile.Username, test.description)
		assert.Equal(t, "admin@testination.com", profile.Email, test.description)
		assert.Equal(t, 100, profile.Levels[0].Score, test.description)
		assert.Equal(t, 0, profile.Levels[0].TimeFreezePointsUsed, test.description)
	}
	/*
		type ProfileData = {
			svg: string;
			username: string;
			email: string;
			levels: {
				title: string;
				score: number;
				time_freeze_points_used: number | null;
				textual_hint_points_used: number | null;
				hint_solution_points_used: number | null;
				start_time: string;
				end_time: string;
			}[];
		};
	*/
}

func TestAvailableIcons(t *testing.T) {
	tests := []struct {
		description string
		username    string
		password    string
	}{
		{
			description: "Get available icons for pic selection",
			username:    "admin",
			password:    "rootroot",
		},
	}

	db := database.CreateFinalTestinationDB(
		env.DB_USERNAME, env.DB_PASSWORD,
		env.DB_HOST, env.DB_PORT, env.DB_NAME)
	db.CreateSchemas()

	app := fiber.New()
	playerGroup := app.Group("/player")
	SetUpPlayerRoutes(&playerGroup, db)

	for _, test := range tests {
		loginResp := utils.MockLogin(t, app, test.username, test.password)
		cookie := utils.Filter(loginResp.Cookies(), func(c *http.Cookie) bool {
			return c.Name == constants.AUTH_COOKIE_NAME
		})[0].Value
		req := httptest.NewRequest("GET", "/player/availableIcons", nil)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: constants.AUTH_COOKIE_NAME, Value: cookie})

		resp, errTest := app.Test(req, -1)
		assert.NoError(t, errTest)
		var availableIcons []functionality.IconDTO
		body, errBody := io.ReadAll(resp.Body)
		assert.NoError(t, errBody)
		errUnmarshal := json.Unmarshal(body, &availableIcons)
		assert.NoError(t, errUnmarshal)
		assert.Equal(t, 200, resp.StatusCode, test.description)
		assert.Equal(t, 9, len(availableIcons), test.description)

		for index, icon := range availableIcons {
			assert.NotEmpty(t, icon.Id, test.description)
			assert.NotEmpty(t, icon.Svg, test.description)
			assert.Equal(t, fmt.Sprintf("%d", index+1), icon.Id, test.description)
			assert.True(t, strings.HasPrefix(icon.Svg, "<svg "), test.description)
			assert.True(t, strings.HasSuffix(icon.Svg, "</svg>"), test.description)
		}
	}
}

var changeIconRoute string = "/player/changeIcon"
var changeIconBody string = `{"icon": "2"}`

func TestChangeIcon(t *testing.T) {
	tests := []struct {
		method       string
		description  string // description of the test case
		username     string
		password     string
		route        string // route path to test
		expectedCode int    // expected HTTP status code
		body         string
		expected     string
	}{
		{
			method:       "POST",
			description:  "Player changes icon",
			username:     "admin",
			password:     "rootroot",
			route:        changeIconRoute,
			expectedCode: 200,
			body:         changeIconBody,
		},
	}

	db := database.CreateFinalTestinationDB(
		env.DB_USERNAME, env.DB_PASSWORD,
		env.DB_HOST, env.DB_PORT, env.DB_NAME)
	db.CreateSchemas()

	app := fiber.New()
	playerGroup := app.Group("/player")
	SetUpPlayerRoutes(&playerGroup, db)

	for _, test := range tests {
		loginResp := utils.MockLogin(t, app, test.username, test.password)
		cookie := utils.Filter(loginResp.Cookies(), func(c *http.Cookie) bool {
			return c.Name == constants.AUTH_COOKIE_NAME
		})[0].Value
		req := httptest.NewRequest(test.method, test.route, bytes.NewBuffer([]byte(test.body)))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: constants.AUTH_COOKIE_NAME, Value: cookie})

		resp, err := app.Test(req, -1) // -1 means no timeout

		assert.NoError(t, err)
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
func TestLogOut(t *testing.T) {
	tests := []struct {
		description  string
		username     string
		password     string
		expectedCode int
	}{
		{
			description:  "Log out",
			username:     "admin",
			password:     "rootroot",
			expectedCode: 200,
		},
	}
	db := database.CreateFinalTestinationDB(
		env.DB_USERNAME, env.DB_PASSWORD,
		env.DB_HOST, env.DB_PORT, env.DB_NAME)
	db.CreateSchemas()

	app := fiber.New()
	playerGroup := app.Group("/player")
	SetUpPlayerRoutes(&playerGroup, db)
	for _, test := range tests {
		loginResp := utils.MockLogin(t, app, test.username, test.password)
		assert.Equal(t, test.expectedCode, loginResp.StatusCode, test.description)
		cookie := utils.Filter(loginResp.Cookies(), func(c *http.Cookie) bool {
			return c.Name == constants.AUTH_COOKIE_NAME
		})
		if test.expectedCode == 200 {
			assert.NotEmpty(t, cookie, test.description)
			assert.NotEmpty(t, cookie[0].Value, test.description)
		}
		req := httptest.NewRequest("POST", "/player/logout", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, errTest := app.Test(req, -1)
		cookie = utils.Filter(resp.Cookies(), func(c *http.Cookie) bool {
			return c.Name == constants.AUTH_COOKIE_NAME

		})
		assert.NoError(t, errTest)
		assert.Equal(t, 200, resp.StatusCode, test.description)
		assert.Empty(t, cookie[0].Value, test.description)

	}
}
