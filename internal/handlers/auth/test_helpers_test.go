package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func jsonRequest(t *testing.T, app *fiber.App, method string, path string, body map[string]string) *http.Response {
	t.Helper()

	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(method, path, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	return resp
}

func signupWithStatus(t *testing.T, app *fiber.App, username string, password string, expectedStatus int) *http.Response {
	t.Helper()

	resp := jsonRequest(t, app, http.MethodPost, "/api/auth/signup", map[string]string{
		"username": username,
		"password": password,
	})
	require.Equal(t, expectedStatus, resp.StatusCode)

	return resp
}

func loginWithStatus(t *testing.T, app *fiber.App, username string, password string, expectedStatus int) *http.Response {
	t.Helper()

	resp := jsonRequest(t, app, http.MethodPost, "/api/auth/login", map[string]string{
		"username": username,
		"password": password,
	})
	require.Equal(t, expectedStatus, resp.StatusCode)

	return resp
}

func refreshWithStatus(t *testing.T, app *fiber.App, refreshToken string, expectedStatus int) *http.Response {
	t.Helper()

	resp := jsonRequest(t, app, http.MethodPost, "/api/auth/refresh", map[string]string{
		"refresh_token": refreshToken,
	})
	require.Equal(t, expectedStatus, resp.StatusCode)

	return resp
}

func decodeTokenResponse(t *testing.T, resp *http.Response) tokenResponse {
	t.Helper()

	var tokens tokenResponse
	err := json.NewDecoder(resp.Body).Decode(&tokens)
	require.NoError(t, err)
	require.NotEmpty(t, tokens.AccessToken)
	require.NotEmpty(t, tokens.RefreshToken)

	return tokens
}
