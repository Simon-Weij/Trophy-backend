package clips

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api")

	Register(api, nil)

	t.Run("get clip route exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/clips/abcdef1234567890", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("upload route is protected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/clips/", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("delete route is protected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/clips/abcdef1234567890", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}
