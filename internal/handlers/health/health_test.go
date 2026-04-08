package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"trophy/internal/dbtest"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	db := dbtest.SetupDB(t)
	handler := NewHandler(db)

	app := fiber.New()
	app.Get("/health", handler.Health)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body map[string]any
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	require.Equal(t, "healthy", body["status"])

	checks, ok := body["checks"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "ok", checks["database"])
}

func TestHealth_UnhealthyWhenDatabaseUnavailable(t *testing.T) {
	db := dbtest.SetupDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	handler := NewHandler(db)
	app := fiber.New()
	app.Get("/health", handler.Health)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, testErr := app.Test(req)
	require.NoError(t, testErr)
	require.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
}

func TestProtected(t *testing.T) {
	handler := NewHandler(nil)
	app := fiber.New()
	app.Get("/health/protected", handler.Protected)

	req := httptest.NewRequest(http.MethodGet, "/health/protected", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
}

