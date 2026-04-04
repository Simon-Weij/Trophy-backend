package health

import (
	"trophy/internal/handlers/health"

	"github.com/gofiber/fiber/v3"
)

func Register(api fiber.Router) {
	api.Get("/health", health.Health)
}
