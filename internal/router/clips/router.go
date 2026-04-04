package clips

import (
	"trophy/internal/handlers/auth"
	"trophy/internal/handlers/clips"

	"github.com/gofiber/fiber/v3"
)

func Register(api fiber.Router) {
	clipGroup := api.Group("/clips")
	clipGroup.Post("/", auth.AuthMiddleware, clips.UploadClip)
	clipGroup.Get("/:hash", clips.GetClip)
	clipGroup.Delete("/:hash", auth.AuthMiddleware, clips.DeleteClip)
}
