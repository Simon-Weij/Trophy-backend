package auth_test

import (
	"net/http"
	"testing"
	"trophy/internal/testutils"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestSignup(t *testing.T) {
	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
	}{
		{
			name:           "Valid signup",
			body:           map[string]string{"username": "myusername", "password": "secret123"},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Non-matching json",
			body:           map[string]string{"email": "test@example.com", "password": "secret123"},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "Duplicate username",
			body:           map[string]string{"username": "duplicateuser", "password": "secret123"},
			expectedStatus: fiber.StatusConflict,
		},
		{
			name:           "Missing username",
			body:           map[string]string{"password": "secret123"},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "Missing password",
			body:           map[string]string{"username": "myusername"},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "Malformed JSON",
			body:           map[string]string{"username": "loginuser", "password": "password123", "invalidField": "field"},
			expectedStatus: fiber.StatusBadRequest,
		},
	}
	app, _ := testutils.SetupApp(t)

	t.Run("Create first user", func(t *testing.T) {
		signupWithStatus(t, app, "duplicateuser", "passwordabc", fiber.StatusOK)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := jsonRequest(t, app, http.MethodPost, "/api/auth/signup", tt.body)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
