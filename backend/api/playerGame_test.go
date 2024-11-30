package api

import (
	"backend/database"
	"backend/env"
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type expectedLeaderBoardResponse struct {
	CurrentPage int `json:"currentPage"`
	Pages       int `json:"pages"`
	Entries     []struct {
		Username string `json:"username"`
		Score    int    `json:"score"`
	} `json:"entries"`
}

// Unit testing for sunction TestMatchArrays
func TestLeaderBoardNavigation(t *testing.T) {
	tests := []struct {
		description           string
		inputPage             string
		expectedPage          int
		expectedMaxPage       int
		expectedEntriesLength int
	}{
		{
			description:           "When I request the first page of the leaderboard, I should get the first page of the leaderboard",
			inputPage:             "1",
			expectedPage:          1,
			expectedMaxPage:       2,
			expectedEntriesLength: 25,
		},
		{
			description:           "When I request a page that exceeds the max page number, I should get the last page of the leaderboard",
			inputPage:             "3",
			expectedPage:          2,
			expectedMaxPage:       2,
			expectedEntriesLength: 18,
		},
		{
			description:           "When I request a page below the first page, I should get the first page of the leaderboard",
			inputPage:             "0",
			expectedPage:          1,
			expectedMaxPage:       2,
			expectedEntriesLength: 25,
		},
		{
			description:           "When I request a page with invalid page number, I should get the first page of the leaderboard",
			inputPage:             "invalid",
			expectedPage:          1,
			expectedMaxPage:       2,
			expectedEntriesLength: 25,
		},
	}

	db := database.CreateFinalTestinationDB(
		env.DB_USERNAME, env.DB_PASSWORD,
		env.DB_HOST, env.DB_PORT, env.DB_NAME)
	db.CreateSchemas()

	app := fiber.New()
	gameGroup := app.Group("/leaderboard")
	SetUpPlayerGameRoutes(&gameGroup, db)

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/leaderboard/"+test.inputPage, bytes.NewBuffer([]byte("")))

		resp, err := app.Test(req, -1) // -1 means no timeout

		assert.NoError(t, err)
		assert.Equalf(t, 200, resp.StatusCode, test.description)

		responseBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var parsedResponseBody expectedLeaderBoardResponse
		err = json.Unmarshal(responseBody, &parsedResponseBody)
		assert.NoError(t, err)

		assert.Equalf(t, test.expectedPage, parsedResponseBody.CurrentPage, test.description)
		assert.Equalf(t, test.expectedMaxPage, parsedResponseBody.Pages, test.description)
	}
}
