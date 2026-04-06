package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"trophy/internal/testutils"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestRefresh(t *testing.T) {
	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
	}{
		{
			name:           "Missing refresh token",
			body:           map[string]string{},
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Invalid refresh token",
			body:           map[string]string{"refresh_token": "not-a-real-token"},
			expectedStatus: fiber.StatusUnauthorized,
		},
	}
	app, _ := testutils.SetupApp(t)

	t.Run("Create first user", func(t *testing.T) {
		signupWithStatus(t, app, "refreshuser", "password", fiber.StatusOK)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := jsonRequest(t, app, http.MethodPost, "/api/auth/refresh", tt.body)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})

	}

	t.Run("Valid refresh rotates tokens", func(t *testing.T) {
		uniqueUsername := "refreshuser-rotation"
		signupResp := signupWithStatus(t, app, uniqueUsername, "password123", fiber.StatusOK)
		firstTokens := decodeTokenResponse(t, signupResp)

		refreshResp := refreshWithStatus(t, app, firstTokens.RefreshToken, fiber.StatusOK)
		rotatedTokens := decodeTokenResponse(t, refreshResp)

		require.NotEqual(t, firstTokens.RefreshToken, rotatedTokens.RefreshToken)
	})

	t.Run("Used refresh token cannot be reused", func(t *testing.T) {
		uniqueUsername := "refreshuser-reuse"
		signupResp := signupWithStatus(t, app, uniqueUsername, "password123", fiber.StatusOK)
		firstTokens := decodeTokenResponse(t, signupResp)

		refreshWithStatus(t, app, firstTokens.RefreshToken, fiber.StatusOK)
		refreshWithStatus(t, app, firstTokens.RefreshToken, fiber.StatusUnauthorized)
	})

	t.Run("Refreshed access token works on protected route", func(t *testing.T) {
		uniqueUsername := "refreshuser-protected"
		signupResp := signupWithStatus(t, app, uniqueUsername, "password123", fiber.StatusOK)
		firstTokens := decodeTokenResponse(t, signupResp)

		refreshResp := refreshWithStatus(t, app, firstTokens.RefreshToken, fiber.StatusOK)
		rotatedTokens := decodeTokenResponse(t, refreshResp)

		protectedReq := httptest.NewRequest(http.MethodGet, "/api/health/protected", nil)
		protectedReq.Header.Set("Authorization", "Bearer "+rotatedTokens.AccessToken)

		protectedResp, err := app.Test(protectedReq)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusOK, protectedResp.StatusCode)
	})
}
