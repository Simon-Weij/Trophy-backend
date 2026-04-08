package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"trophy/internal/dbtest"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	app := fiber.New()
	db := dbtest.SetupDB(t)
	api := app.Group("/api")

	Register(api, db)

	tests := []struct {
		name string
		path string
	}{
		{name: "login route exists", path: "/api/auth/login"},
		{name: "signup route exists", path: "/api/auth/signup"},
		{name: "refresh route exists", path: "/api/auth/refresh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusMethodNotAllowed, resp.StatusCode)
		})
	}
}
