package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestPingServer(t *testing.T) {
	tests := []struct {
		description  string // description of the test case
		route        string // route path to test
		expectedCode int    // expected HTTP status code
	}{
		{
			description:  "Ping server",
			route:        "/ping",
			expectedCode: 200,
		},
	}

	app := fiber.New()
	app.Get("/ping", Ping)

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, nil)
		resp, _ := app.Test(req, -1) // -1 means no timeout
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
