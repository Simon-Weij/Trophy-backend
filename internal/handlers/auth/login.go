package auth

import (
	"trophy/internal/database"
	apphttp "trophy/internal/http"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login authenticates a user and returns a token pair.
// @Summary User login
// @Description Sign in with username and password to get access tokens
// @Tags Auth
// @Param request body loginRequest true "Login credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401
// @Failure 500
// @Router /auth/login [post]
func Login(c fiber.Ctx) error {
	var body loginRequest
	if err := apphttp.Bind(c, &body); err != nil {
		return err
	}

	var user database.User
	if err := database.DB.First(&user, "username = ?", body.Username).Error; err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	tokens, err := GenerateTokenPair(user)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(tokens)
}
