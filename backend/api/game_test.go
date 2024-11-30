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
	"slices"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var existingGameResponseBodyAsString gameResponse = gameResponse{
	Blocks: []string{
		"goodcompany",
		"<img",
		"<iframe",
		"50%",
		"evilcompany",
		"<script",
		"25%",
		"/>",
		"></script",
		"></iframe>",
		"100%",
	},
	Cheatsheet: "1. Introduction\nCross-site scripting (XSS) attacks target vulnerable web applications, injecting client-side scripts into web pages viewed by other users. A common way to target someone is to send a link to a webpage that contains\nthe XSS attack so that the target will open the link and inject the script. They bypass the Same Origin Policy (SOP) and may bypass access control systems.\n\n2. Targets and Effects\nThe injection works client-side and so the main target are users, be that directly by sending a malicious link (like in a scam mail) or by storing it in a vulnerable webpage that is going to be loaded by another user.\nXSS attacks can lead to stolen credentials, disclosure of the user's session cookies installation of Trojans or alter the content of a web page that can lead to disastrous consequences.\n\n3. Example\nXSS attacks can be used to steal crediants by tricking users into submitting their data in page they think is the real one when in reality they are writing in a malicious page loaded by the script inside the URL.\nAn XSS attack on a pharmaceutical site can alter the dosage shown to the user, possibly causing an overdose; stealing a session's cookies lets the attacker hijack the user's session and take over the account. These\nare but a few of the possible attacks.\n\n4. Countermeasures\nSanitization of the code is the main way to defend a web application from this kind of attacks. Secure pages should escape all characters that can be used for such an attack (\\\n', \", >, and so on) and write the code\nkeeping in mind the possible attacks that can be used: a sanitized page written poorly can still be attacked in certain scenarios. There are automatic functions that can handle this, like escapeshellcmd(),\nhtmlspecialchars() and escapeshellarg(). There are other ways to mitigate this attacks, like using cookies to handle authentication or by disabling scripts entirely but they are not perfect solutions or can reduce\nthe functionality and responsiveness of the web page.\n\n5. Useful links\nWikipedia",
	Skeleton: map[string]string{
		"1": "src=\"http://",
		"3": ".com\" style=\"position: absolute; top: 0; left: 0; width: ",
		"5": "; height: 100%;\"",
	},
	SolutionLength: 7,
	Story:          "====EMAIL FROM WORK====\nFrom: [MarkusPumpkin@greatTesters.com]\nTo: [007@greatTesters.com]\n\nSubject: Operation Final Testination\nGreetings,\n\nI write to you in the most dire of situations: Q Division has been attacked, probably by a double agent. The whole team got poisoned and those who are not dead are still recovering.  must know of \nOperation Final Testination and tried to sabotage our efforts to break into their system. I cannot help you with all of this stuff: I barely understand what they were talking about, so I will just send whatever data I think\nmight help you in your mission. Apparently they were working on some kind of attack to steal Philip Rich's credentials to get full access to his account on Black Millstone: I'll attach the project directories to this email.\nWhen I was walking by yesterday, I overheard they were planning some sort of \"frame attack\"? Something that should substitute Black Millstone's login page with their own to trick him into typing his login, or something like that.\nDon't know much about it, so you are on your own.\n\nGood luck,\nMarkus Pumpkin\n\nAttachments:\n<evil-company-iframe-login.zip>",
	Title:          "XSS Attack",
	WinningMessage: "Great! You got the credentials: username = fishM4$t3r password = myfishisbigLOL123",
}

type gameResponse struct {
	Title          string            `json:"title"`
	Story          string            `json:"story"`
	Cheatsheet     string            `json:"cheatsheet"`
	SolutionLength int               `json:"solution_length"`
	Blocks         []string          `json:"blocks"`
	Skeleton       map[string]string `json:"skeleton"`
	Background     string            `json:"background"`
	WinningMessage string            `json:"winning_message"`
}

type expectedGameResponse struct {
	gameResponse

	Background string `json:"background"`
}

type expectedHintResponse struct {
	HintContent string `json:"hintContent"`
}

func TestGetGame(t *testing.T) {
	tests := []struct {
		method       string
		description  string // description of the test case
		route        string // route path to test
		expectedCode int    // expected HTTP status code
		body         string
		expected     gameResponse
	}{
		{
			method:       "GET",
			description:  "Get a game that exists",
			route:        "/game/af8e4754-1b84-4fec-bec4-154a3f894b8f",
			expectedCode: 200,
			expected:     existingGameResponseBodyAsString,
		},
		{
			method:       "GET",
			description:  "Get a game that does not exist",
			route:        "/game/0987afd7-474b-4308-9f2f-447a0995a1ae",
			expectedCode: 404,
		},
		{
			method:       "GET",
			description:  "Get a game with an invalid UUID",
			route:        "/game/af8e4754-1b84-4fec-bec4-154a3f894b811",
			expectedCode: 400,
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
		loginResp := utils.MockLogin(t, app, "admin", "rootroot")
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

		var parsedResponseBody expectedGameResponse
		err = json.Unmarshal(responseBody, &parsedResponseBody)
		assert.NoError(t, err)

		assert.Equalf(t, test.expected.Title, parsedResponseBody.Title, test.description)
		assert.Equalf(t, test.expected.Story, parsedResponseBody.Story, test.description)
		assert.Equalf(t, test.expected.Cheatsheet, parsedResponseBody.Cheatsheet, test.description)
		assert.Equalf(t, test.expected.SolutionLength, parsedResponseBody.SolutionLength, test.description)

		slices.Sort(test.expected.Blocks)
		slices.Sort(parsedResponseBody.Blocks)
		assert.Equalf(t, test.expected.Blocks, parsedResponseBody.Blocks, test.description)

		assert.Equalf(t, test.expected.Skeleton, parsedResponseBody.Skeleton, test.description)
		assert.Equalf(t, test.expected.WinningMessage, parsedResponseBody.WinningMessage, test.description)

		assert.True(t, strings.Contains(parsedResponseBody.Background, "svg"), test.description)
	}
}

func TestUseHint(t *testing.T) {
	tests := []struct {
		method               string
		description          string
		route                string
		expectedCode         int
		body                 string
		expectedHintResponse string
	}{
		{
			method:               "POST",
			description:          "Use a hint",
			route:                "/game/05732286-9fa5-45d4-bef3-13ae0d481afa/hint",
			expectedCode:         200,
			body:                 `{"hint_type":"textual"}`,
			expectedHintResponse: `{"hintContent":"We are trying to get the last ticket purchased by the username we got from the first level so that we can have details about it. We need to close the quotes in the query to write the attack."}`,
		},
		{
			method:               "POST",
			description:          "Use an invalid hint",
			route:                "/game/af8e4754-1b84-4fec-bec4-154a3f894b8f/hint",
			expectedCode:         400,
			body:                 `{"hint_type":"invalid"}`,
			expectedHintResponse: `{"error":"Invalid hint type"}`,
		},
		{
			method:               "POST",
			description:          "Use a hint a second time",
			route:                "/game/05732286-9fa5-45d4-bef3-13ae0d481afa/hint",
			expectedCode:         403,
			body:                 `{"hint_type":"textual"}`,
			expectedHintResponse: `{"error":"cannot buy hint with type textual a second time"}`,
		},
		{
			method:               "POST",
			description:          "Use a hint with not enough coins",
			route:                "/game/05732286-9fa5-45d4-bef3-13ae0d481afa/hint",
			expectedCode:         400,
			body:                 `{"hint_type":"freeze"}`,
			expectedHintResponse: `{"error":"not enough coins"}`,
		},
		{
			method:               "POST",
			description:          "Use fill hint with invalid order",
			route:                "/game/af8e4754-1b84-4fec-bec4-154a3f894b8f/hint",
			expectedCode:         400,
			body:                 `{"hint_type":"fill","order":10}`,
			expectedHintResponse: `{"error":"block not found"}`,
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
		user := "test"
		if test.description == "Use a hint with not enough coins" {
			user = "PleaseRunTests"
		}
		loginResp := utils.MockLogin(t, app, user, "rootroot")
		cookie := utils.Filter(loginResp.Cookies(), func(c *http.Cookie) bool {
			return c.Name == constants.AUTH_COOKIE_NAME
		})[0].Value

		pre_req := httptest.NewRequest("GET", strings.TrimSuffix(test.route, "/hint"), nil)
		pre_req.Header.Set("Content-Type", "application/json")
		pre_req.AddCookie(&http.Cookie{Name: constants.AUTH_COOKIE_NAME, Value: cookie})
		_, err := app.Test(pre_req, -1) // -1 means no timeout
		if err != nil {
			t.Errorf("Error making pre request: %v", err)
		}

		req := httptest.NewRequest(test.method, test.route, bytes.NewBuffer([]byte(test.body)))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: constants.AUTH_COOKIE_NAME, Value: cookie})

		resp, err := app.Test(req, -1) // -1 means no timeout

		assert.NoError(t, err)
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		responseBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			if test.description == "Use a hint" {
				var parsedResponseBody expectedHintResponse
				err = json.Unmarshal(responseBody, &parsedResponseBody)
				assert.NoError(t, err)

				var expectedParsedResponseBody expectedHintResponse
				err = json.Unmarshal([]byte(test.expectedHintResponse), &expectedParsedResponseBody)
				assert.NoError(t, err)
				assert.Equalf(t, expectedParsedResponseBody, parsedResponseBody, test.description)
			}
		} else {
			var parsedResponseBody errorResponse
			err = json.Unmarshal(responseBody, &parsedResponseBody)
			assert.NoError(t, err)
			var expectedErrorParsedResponseBody errorResponse
			err = json.Unmarshal([]byte(test.expectedHintResponse), &expectedErrorParsedResponseBody)
			assert.NoError(t, err)
			assert.Equalf(t, expectedErrorParsedResponseBody, parsedResponseBody, test.description)
		}
	}
}
