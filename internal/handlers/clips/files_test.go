package clips

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestGetClip(t *testing.T) {
	oldBaseDir := clipsBaseDir
	t.Cleanup(func() {
		clipsBaseDir = oldBaseDir
	})

	baseDir := t.TempDir()
	clipsBaseDir = baseDir

	hash := "abcdef1234567890"
	clipDir := filepath.Join(baseDir, hash[0:2], hash[2:4], hash[4:6])
	require.NoError(t, os.MkdirAll(clipDir, 0755))

	clipPath := filepath.Join(clipDir, hash[6:]+".webm")
	require.NoError(t, os.WriteFile(clipPath, []byte("video-data"), 0600))

	app := fiber.New()
	app.Get("/clips/:hash", GetClip)

	req := httptest.NewRequest(http.MethodGet, "/clips/"+hash, nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "video-data", string(bodyBytes))
}

func TestGetClip_NotFound(t *testing.T) {
	oldBaseDir := clipsBaseDir
	t.Cleanup(func() {
		clipsBaseDir = oldBaseDir
	})

	clipsBaseDir = t.TempDir()
	hash := "abcdef1234567890"

	app := fiber.New()
	app.Get("/clips/:hash", GetClip)

	req := httptest.NewRequest(http.MethodGet, "/clips/"+hash, nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestDeleteFile(t *testing.T) {
	oldBaseDir := clipsBaseDir
	t.Cleanup(func() {
		clipsBaseDir = oldBaseDir
	})

	baseDir := t.TempDir()
	clipsBaseDir = baseDir

	hash := "abcdef1234567890"
	clipDir := filepath.Join(baseDir, hash[0:2], hash[2:4], hash[4:6])
	require.NoError(t, os.MkdirAll(clipDir, 0755))

	clipPath := filepath.Join(clipDir, hash[6:]+".webm")
	require.NoError(t, os.WriteFile(clipPath, []byte("video-data"), 0600))

	require.NoError(t, deleteFile(hash))

	_, err := os.Stat(clipPath)
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))

	require.NoError(t, deleteFile(hash))
}
