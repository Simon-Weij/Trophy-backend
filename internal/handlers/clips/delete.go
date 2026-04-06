package clips

import (
	"errors"
	"trophy/internal/database"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"gorm.io/gorm"
)

// DeleteClip deletes a clip owned by the authenticated user.
// @Summary Delete your clip
// @Description Delete a clip you own by its hash
// @Tags Clips
// @Param hash path string true "Clip hash"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clips/{hash} [delete]
func (handler *Handler) DeleteClip(c fiber.Ctx) error {
	userID, err := handler.getUserID(c)
	if err != nil {
		log.Errorf("Couldn't get userID: %v", err)
		return err
	}

	hash := c.Params("hash")

	userOwnsClip, err := handler.userOwnsClip(userID, hash)
	if err != nil {
		log.Errorf("Couldn't check if the user owned the clip for %s", userID)
		return fiber.ErrInternalServerError
	}

	if !userOwnsClip {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "you don't own this video",
		})
	}

	var clip database.Clip
	result := handler.db.Delete(&clip, "user_id = ? AND video_hash = ?", userID, hash)
	if result.Error != nil {
		log.Errorf("Couldn't delete clip for user %d and hash %s: %v", userID, hash, result.Error)
		return fiber.ErrInternalServerError
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "video doesn't exist",
		})
	}
	if err := deleteFile(hash); err != nil {
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(fiber.StatusOK)
}

func deleteFile(hash string) error {
	filePath := filepath.Join(
		clipsBaseDir,
		hash[0:2],
		hash[2:4],
		hash[4:6],
		hash[6:]+".webm",
	)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func (handler *Handler) userOwnsClip(userID uint, hash string) (bool, error) {
	var clip database.Clip
	result := handler.db.Where("user_id = ? AND video_hash = ?", userID, hash).First(&clip)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil

}
