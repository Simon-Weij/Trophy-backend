package auth

import (
	"trophy/internal/handlers/auth"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Register(api fiber.Router, db *gorm.DB) {
	h := auth.NewHandler(db)

	authGroup := api.Group("/auth")
	authGroup.Post("/login", h.Login)
	authGroup.Post("/signup", h.Signup)
	authGroup.Post("/refresh", h.Refresh)
}
