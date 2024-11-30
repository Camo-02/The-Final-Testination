package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func MockLogin(t *testing.T, app *fiber.App, credential string, password string) *http.Response {
	credentials := map[string]string{
		"credential": credential,
		"password":   password,
	}

	jsonCredentials, _ := json.Marshal(credentials)

	postReq := httptest.NewRequest("POST", "/player/login", bytes.NewBuffer(jsonCredentials))
	postReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(postReq, -1) // -1 means no timeout

	assert.NoError(t, err)
	return resp
}
