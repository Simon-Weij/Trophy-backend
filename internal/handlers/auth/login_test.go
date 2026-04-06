package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"trophy/internal/testutils"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
	}{
		{
			name:           "Valid login",
			body:           map[string]string{"username": "loginuser", "password": "password123"},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Wrong password",
			body:           map[string]string{"username": "loginuser", "password": "notthepassword"},
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Non-existent username",
			body:           map[string]string{"username": "nonexistentuser", "password": "anypassword"},
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Empty request body",
			body:           map[string]string{},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "Malformed JSON",
			body:           map[string]string{"username": "loginuser", "password": "password123", "invalidField": "field"},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "No username",
			body:           map[string]string{"password": "password123"},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "No password",
			body:           map[string]string{"username": "loginuser"},
			expectedStatus: fiber.StatusBadRequest,
		},
	}
	app, _ := testutils.SetupApp(t)

	t.Run("Create first user", func(t *testing.T) {
		signupWithStatus(t, app, "loginuser", "password123", fiber.StatusOK)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := jsonRequest(t, app, http.MethodPost, "/api/auth/login", tt.body)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
	t.Run("Should return 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health/protected", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
	t.Run("Should return 200", func(t *testing.T) {
		resp := loginWithStatus(t, app, "loginuser", "password123", fiber.StatusOK)
		tokens := decodeTokenResponse(t, resp)
		accessToken := tokens.AccessToken
		require.NotEmpty(t, accessToken)

		protectedReq := httptest.NewRequest(http.MethodGet, "/api/health/protected", nil)
		protectedReq.Header.Set("Authorization", "Bearer "+accessToken)

		protectedResp, err := app.Test(protectedReq)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusOK, protectedResp.StatusCode)
	})
}
