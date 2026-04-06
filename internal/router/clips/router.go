package clips

import (
	"trophy/internal/handlers/auth"
	"trophy/internal/handlers/clips"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Register(api fiber.Router, db *gorm.DB) {
	h := clips.NewHandler(db)

	clipGroup := api.Group("/clips")
	clipGroup.Post("/", auth.AuthMiddleware, h.UploadClip)
	clipGroup.Get("/:hash", clips.GetClip)
	clipGroup.Delete("/:hash", auth.AuthMiddleware, h.DeleteClip)
}
