package auth

import (
	"trophy/internal/handlers/auth"

	"github.com/gofiber/fiber/v3"
)

func Register(api fiber.Router) {
	authGroup := api.Group("/auth")
	authGroup.Post("/login", auth.Login)
	authGroup.Post("/signup", auth.Signup)
	authGroup.Post("/refresh", auth.Refresh)
}
