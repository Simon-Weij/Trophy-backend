package health

import (
	"trophy/internal/handlers/auth"
	"trophy/internal/handlers/health"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Register(api fiber.Router, db *gorm.DB) {
	h := health.NewHandler(db)
	api.Get("/health", h.Health)
	api.Get("/health/protected", auth.AuthMiddleware, h.Protected)
}
