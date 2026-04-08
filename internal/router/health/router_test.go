package health

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

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	protectedReq := httptest.NewRequest(http.MethodGet, "/api/health/protected", nil)
	protectedResp, testErr := app.Test(protectedReq)
	require.NoError(t, testErr)
	require.Equal(t, fiber.StatusUnauthorized, protectedResp.StatusCode)
}
