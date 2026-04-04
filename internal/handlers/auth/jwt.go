package auth

import (
	"trophy/internal/database"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func jwtSigningKey() []byte {
	return []byte(os.Getenv("JWT_KEY"))
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshHandler(c fiber.Ctx) error {
	var body RefreshRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var storedToken database.RefreshToken
	if err := database.DB.First(&storedToken, "token = ? AND expires_at > ?", body.RefreshToken, time.Now()).Error; err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var user database.User
	if err := database.DB.First(&user, storedToken.UserID).Error; err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	database.DB.Delete(&storedToken)

	newTokens, err := GenerateTokenPair(user)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(newTokens)
}

func AuthMiddleware(c fiber.Ctx) error {
	authHeader := strings.TrimSpace(c.Get("Authorization"))
	if authHeader == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	tokenStr := authHeader
	if len(authHeader) > len("Bearer ") && strings.EqualFold(authHeader[:len("Bearer")], "Bearer") {
		tokenStr = strings.TrimSpace(authHeader[len("Bearer"):])
	}

	if tokenStr == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSigningKey(), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("username", claims.Username)

	return c.Next()
}
