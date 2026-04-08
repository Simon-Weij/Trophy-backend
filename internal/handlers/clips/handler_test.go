package clips

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"trophy/internal/database"
	"trophy/internal/dbtest"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	db := dbtest.SetupDB(t)
	handler := NewHandler(db)
	require.NotNil(t, handler)
}

func TestGetUserID(t *testing.T) {
	db := dbtest.SetupDB(t)
	handler := NewHandler(db)

	app := fiber.New()
	app.Get("/user-id", func(c fiber.Ctx) error {
		c.Locals("username", "clips-user")
		_, err := handler.getUserID(c)
		require.NoError(t, err)
		return c.SendStatus(fiber.StatusNoContent)
	})

	require.NoError(t, db.Create(&database.User{Username: "clips-user", Password: "password"}).Error)

	req := httptest.NewRequest(http.MethodGet, "/user-id", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestGetUserID_UnauthorizedWhenMissingUsername(t *testing.T) {
	handler := NewHandler(nil)

	app := fiber.New()
	app.Get("/user-id", func(c fiber.Ctx) error {
		_, err := handler.getUserID(c)
		require.ErrorIs(t, err, fiber.ErrUnauthorized)
		return c.SendStatus(fiber.StatusUnauthorized)
	})

	req := httptest.NewRequest(http.MethodGet, "/user-id", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestDeleteClip(t *testing.T) {
	oldBaseDir := clipsBaseDir
	t.Cleanup(func() {
		clipsBaseDir = oldBaseDir
	})

	baseDir := t.TempDir()
	clipsBaseDir = baseDir

	db := dbtest.SetupDB(t)
	handler := NewHandler(db)

	user := database.User{Username: "owner", Password: "password"}
	require.NoError(t, db.Create(&user).Error)

	hash := "abcdef1234567890"
	clip := database.Clip{Title: "owned clip", VideoHash: hash, UserID: user.ID}
	require.NoError(t, db.Create(&clip).Error)

	clipDir := filepath.Join(baseDir, hash[0:2], hash[2:4], hash[4:6])
	require.NoError(t, os.MkdirAll(clipDir, 0755))
	clipPath := filepath.Join(clipDir, hash[6:]+".webm")
	require.NoError(t, os.WriteFile(clipPath, []byte("video-data"), 0600))

	app := fiber.New()
	app.Delete("/clips/:hash", func(c fiber.Ctx) error {
		c.Locals("username", "owner")
		return handler.DeleteClip(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/clips/"+hash, nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	_, statErr := os.Stat(clipPath)
	require.Error(t, statErr)
	require.True(t, os.IsNotExist(statErr))
}

func TestUploadClip_MissingVideo(t *testing.T) {
	db := dbtest.SetupDB(t)
	handler := NewHandler(db)

	app := fiber.New()
	app.Post("/clips", func(c fiber.Ctx) error {
		return handler.UploadClip(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/clips", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
