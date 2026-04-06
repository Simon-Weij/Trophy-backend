package clips

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"trophy/internal/database"
	apphttp "trophy/internal/http"
	"trophy/internal/video"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"gorm.io/gorm"
)

type uploadRequest struct {
	Title string `json:"title" form:"title"`
}

var clipsBaseDir = "/var/uploads"

type clipActionResponse struct {
	Message string `json:"message"`
	Hash    string `json:"hash"`
}

// UploadClip uploads a new clip and stores metadata.
// @Summary Upload a clip
// @Description Upload a video clip (automatically converted to WebM if needed)
// @Tags Clips
// @Accept multipart/form-data
// @Param title formData string false "Optional title for your clip"
// @Param video formData file true "Video file to upload"
// @Security BearerAuth
// @Success 200 {object} clipActionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clips/ [post]
func (handler *Handler) UploadClip(c fiber.Ctx) error {
	var body uploadRequest
	if err := apphttp.Bind(c, &body); err != nil {
		return err
	}

	file, err := c.FormFile("video")
	if err != nil {
		return fiber.ErrBadRequest
	}

	hash, err := hashFile(file)
	if err != nil {
		log.Errorf("Hash failed for file %s: %v", file.Filename, err)
		return fiber.ErrInternalServerError
	}

	var existingClip database.Clip
	if err := handler.db.Select("id").Where("video_hash = ?", hash).First(&existingClip).Error; err == nil {
		return c.JSON(fiber.Map{
			"message": "Video already exists",
			"hash":    hash,
		})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("Failed to check clip existence for hash %s: %v", hash, err)
		return fiber.ErrInternalServerError
	}

	if err := saveFile(c, file, hash); err != nil {
		log.Errorf("Failed to store clip %s: %v", hash, err)
		return fiber.ErrInternalServerError
	}

	userID, err := handler.getUserID(c)
	if err != nil {
		return err
	}

	newClip := database.Clip{
		Title:     body.Title,
		VideoHash: hash,
		UserID:    userID,
	}

	if err := handler.db.Create(&newClip).Error; err != nil {
		log.Errorf("Failed to save clip metadata for hash %s: %v", hash, err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"message": "Upload successful",
		"hash":    hash,
	})
}

func (handler *Handler) getUserID(context fiber.Ctx) (uint, error) {
	username, ok := context.Locals("username").(string)
	if !ok || strings.TrimSpace(username) == "" {
		return 0, fiber.ErrUnauthorized
	}

	var user database.User
	if err := handler.db.Select("id").Where("username = ?", username).First(&user).Error; err != nil {
		return 0, fiber.ErrUnauthorized
	}

	return user.ID, nil
}

func hashFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, src); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func saveFile(context fiber.Ctx, file *multipart.FileHeader, hash string) error {
	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("upload_%s_%d.tmp", hash, time.Now().UnixNano()))
	if err := context.SaveFile(file, tempPath); err != nil {
		return err
	}
	defer os.Remove(tempPath)

	return processVideo(tempPath, hash)
}

func processVideo(tempPath string, hash string) error {
	container, _ := video.GetContainer(tempPath)
	codec, _ := video.GetCodec(tempPath)

	uploadDir := filepath.Join(clipsBaseDir, hash[0:2], hash[2:4], hash[4:6])
	uploadPath := filepath.Join(uploadDir, hash[6:]+".webm")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return err
	}

	if !strings.Contains(strings.ToLower(container), "webm") || codec != "vp9" {
		return video.TranscodeToWebm(tempPath, uploadPath)
	}

	return os.Rename(tempPath, uploadPath)
}
