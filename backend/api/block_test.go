package api

import (
	"backend/constants"
	"backend/database"
	"backend/env"
	"backend/utils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Unit testing for sunction TestMatchArrays
func TestMatchArrays(t *testing.T) {
	tests := []struct {
		description     string // description of the test case
		answer          []string
		solution        []string
		expectedMatches []bool
		expectedResult  bool
	}{
		{
			description: "Non-empty arrays correctly match",
			answer: []string{
				"a",
				"b",
				"c",
			},
			solution: []string{
				"a",
				"b",
				"c",
			},
			expectedMatches: []bool{
				true,
				true,
				true,
			},
			expectedResult: true,
		},
		{
			description:     "Empty arrays correctly match",
			answer:          []string{},
			solution:        []string{},
			expectedMatches: []bool{},
			expectedResult:  true,
		},
		{
			description: "Arrays of different lengths do not match",
			answer: []string{
				"a",
				"b",
			},
			solution: []string{
				"a",
				"b",
				"c",
			},
			expectedMatches: []bool{
				false,
				false,
				false,
			},
			expectedResult: false,
		},
		{
			description: "Arrays of same length but different elements do not match",
			answer: []string{
				"a",
				"b",
				"1",
			},
			solution: []string{
				"a",
				"b",
				"c",
			},
			expectedMatches: []bool{
				true,
				true,
				false,
			},
			expectedResult: false,
		},
	}

	for _, test := range tests {
		matches, result := ArraysMatch(test.solution, test.answer)

		assert.Equalf(t, test.expectedMatches, matches, test.description)
		assert.Equalf(t, test.expectedResult, result, test.description)
	}
}

var getLevel string = "/game/af8e4754-1b84-4fec-bec4-154a3f894b8f"

var checkEndpoint string = "/blocks/af8e4754-1b84-4fec-bec4-154a3f894b8f/check-answer"
var correct_submissionString string = `{"blocks": [
    "<iframe",
    "src=\"http://",
    "evilcompany",
    ".com\" style=\"position: absolute; top: 0; left: 0; width: ",
    "100%",
    "; height: 100%;\"",
    "></iframe>"
]}`
var incorrect_submissionString string = `{"blocks": [
    "<iframe",
    "src=\"http://",
    "evilcompany",
    ".com\" style=\"position: absolute; top: 0; left: 0; width: ",
    "50%",
    "; height: 100%;\"",
    "></iframe>"
]}`

var invalid_submissionString string = `{"blocks": "50%"}`
var expected_indexes string = `{"matches":[true, true, true, true, false, true, true]}`
var expected_correct_answer string = `{"next_level_id": "05732286-9fa5-45d4-bef3-13ae0d481afa", "answer_correctly":"correct answer"}`
var expected_error string = `{"error": "Please provide a body that matches the expected structure"}`
var firstLevelEndPoint string = "/blocks/af8e4754-1b84-4fec-bec4-154a3f894b8f/check-answer"
var overwrite_error string = `{"next_level_id": "05732286-9fa5-45d4-bef3-13ae0d481afa"}`

type errorResponse struct {
	Error string `json:"error"`
}

type correctResponseRepeated struct {
	NextLevelID string `json:"next_level_id"`
}

// TODO: can we test score? It depends also on time..
type correctResponse struct {
	AnswerCorrectly string `json:"answer_correctly"`
	NextLevelID     string `json:"next_level_id"`
}

type incorrectResponse struct {
	Matches []bool `json:"matches"`
}

func TestCheckAnswer(t *testing.T) {
	tests := []struct {
		description  string // description of the test case
		endpoint     string
		expectedCode int // expected HTTP status code
		body         string
		expected     string
	}{
		{
			description:  "Submit a wrong solution",
			endpoint:     checkEndpoint,
			expectedCode: 400,
			body:         incorrect_submissionString,
			expected:     expected_indexes,
		},
		{
			description:  "Submit the correct solution",
			endpoint:     checkEndpoint,
			expectedCode: 200,
			body:         correct_submissionString,
			expected:     expected_correct_answer,
		},
		{
			description:  "Submit an invalid solution",
			endpoint:     checkEndpoint,
			expectedCode: 400,
			body:         invalid_submissionString,
			expected:     expected_error,
		},
		{
			description:  "Submit a solution for the second time",
			endpoint:     firstLevelEndPoint,
			expectedCode: 200,
			body:         correct_submissionString,
			expected:     overwrite_error,
		},
	}

	db := database.CreateFinalTestinationDB(
		env.DB_USERNAME, env.DB_PASSWORD,
		env.DB_HOST, env.DB_PORT, env.DB_NAME)
	db.CreateSchemas()

	app := fiber.New()
	blocksGroup := app.Group("/blocks")
	SetUpBlocksRoutes(&blocksGroup, db)

	playerGroup := app.Group("/player")
	SetUpPlayerRoutes(&playerGroup, db)

	gameGroup := app.Group("/game")
	SetUpGameRoutes(&gameGroup, db)

	for _, test := range tests {
		loginResp := utils.MockLogin(t, app, "test", "rootroot")
		cookie := utils.Filter(loginResp.Cookies(), func(c *http.Cookie) bool {
			return c.Name == constants.AUTH_COOKIE_NAME
		})[0].Value
		// get the level to create the db entry
		getReq := httptest.NewRequest("GET", getLevel, bytes.NewBuffer([]byte("")))
		getReq.Header.Set("Content-Type", "application/json")
		getReq.AddCookie(&http.Cookie{Name: constants.AUTH_COOKIE_NAME, Value: cookie})
		_, get_err := app.Test(getReq, -1) // -1 means no timeout
		assert.NoError(t, get_err)
		// post test submission
		postReq := httptest.NewRequest("POST", test.endpoint, bytes.NewBuffer([]byte(test.body)))
		postReq.Header.Set("Content-Type", "application/json")
		postReq.AddCookie(&http.Cookie{Name: constants.AUTH_COOKIE_NAME, Value: cookie})
		resp, err := app.Test(postReq, -1) // -1 means no timeout
		assert.NoError(t, err)
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		responseBody, err := io.ReadAll(resp.Body)

		assert.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			if test.description == "Submit a solution for the second time" {
				var parsedResponseBody correctResponseRepeated
				err = json.Unmarshal(responseBody, &parsedResponseBody)
				assert.NoError(t, err)

				var expectedResponseBody correctResponseRepeated
				err = json.Unmarshal([]byte(test.expected), &expectedResponseBody)
				assert.NoError(t, err)

				assert.Equalf(t, expectedResponseBody, parsedResponseBody, test.description)
			} else if test.description == "Submit the correct solution" {
				var expectedResponseBody correctResponse
				err = json.Unmarshal([]byte(test.expected), &expectedResponseBody)
				assert.NoError(t, err)

				var parsedResponseBody correctResponse
				err = json.Unmarshal(responseBody, &parsedResponseBody)

				//log.Println("\n\n\n", test.description, parsedResponseBody)
				assert.NoError(t, err)

				assert.Equalf(t, expectedResponseBody, parsedResponseBody, test.description)
			}

		} else if resp.StatusCode == 400 {

			if test.description == "Submit a wrong solution" {
				var parsedResponseBody incorrectResponse
				err = json.Unmarshal(responseBody, &parsedResponseBody)
				assert.NoError(t, err)

				var expectedResponseBody incorrectResponse
				err = json.Unmarshal([]byte(test.expected), &expectedResponseBody)
				assert.NoError(t, err)

				assert.Equalf(t, expectedResponseBody, parsedResponseBody, test.description)
			} else if test.description == "Submit an invalid solution" {
				var parsedResponseBody errorResponse
				err = json.Unmarshal(responseBody, &parsedResponseBody)
				assert.NoError(t, err)

				var expectedResponseBody errorResponse
				err = json.Unmarshal([]byte(test.expected), &expectedResponseBody)
				assert.NoError(t, err)

				assert.Equalf(t, expectedResponseBody, parsedResponseBody, test.description)
			}
		}
	}
}
