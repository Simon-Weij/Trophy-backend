package clips

import (
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

// GetClip serves a clip file by its hash.
//
//	@Summary		Download a clip
//	@Description	Get a clip by its hash as a WebM video file
//	@Tags			Clips
//	@Param			hash	path		string	true	"Clip hash"
//	@Success		200		{file}		file
//	@Failure		404		{object}	map[string]string
//	@Router			/clips/{hash} [get]
func GetClip(c fiber.Ctx) error {
	hash := c.Params("hash")

	filePath := filepath.Join(
		clipsBaseDir,
		hash[0:2],
		hash[2:4],
		hash[4:6],
		hash[6:]+".webm",
	)

	log.Infof("Serving clip from path: %s", filePath)

	if err := c.SendFile(filePath); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Clip not found"})
	}

	return nil
}
