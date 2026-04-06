package auth

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func jwtSigningKey() []byte {
	return []byte(os.Getenv("JWT_KEY"))
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
