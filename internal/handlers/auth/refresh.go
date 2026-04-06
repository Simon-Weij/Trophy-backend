package auth

import (
	"time"
	"trophy/internal/database"
	apphttp "trophy/internal/http"

	"github.com/gofiber/fiber/v3"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Refresh rotates a refresh token and returns a new token pair.
// @Summary Get new tokens
// @Description Use a refresh token to get a new access token
// @Tags Auth
// @Param request body refreshRequest true "Refresh token"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401
// @Failure 500
// @Router /auth/refresh [post]
func (handler *Handler) Refresh(c fiber.Ctx) error {
	var body refreshRequest
	if err := apphttp.Bind(c, &body); err != nil {
		return err
	}

	var storedToken database.RefreshToken
	if err := handler.db.First(&storedToken, "token = ? AND expires_at > ?", body.RefreshToken, time.Now()).Error; err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var user database.User
	if err := handler.db.First(&user, storedToken.UserID).Error; err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	handler.db.Delete(&storedToken)

	tokens, err := GenerateTokenPair(handler.db, user)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(tokens)
}
