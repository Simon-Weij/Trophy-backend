package auth

import (
	"errors"
	"trophy/internal/database"
	apphttp "trophy/internal/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type signupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Signup creates a new user account and returns a token pair.
//
//	@Summary		Create new account
//	@Description	Sign up with username and password to create a new account
//	@Tags			Auth
//	@Param			request	body		signupRequest	true	"New account details"
//	@Success		200		{object}	TokenResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/auth/signup [post]
func (handler *Handler) Signup(c fiber.Ctx) error {
	var body signupRequest

	if err := apphttp.Bind(c, &body); err != nil {
		return err
	}

	if body.Username == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user := database.User{
		Username: body.Username,
		Password: string(hashedPassword),
	}

	if err := handler.db.Create(&user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username is already taken",
			})
		}

		log.Errorf("Signup error for %s: %s", body.Username, err)
		return c.Status(fiber.StatusInternalServerError).JSON("Something went wrong")
	}

	tokens, err := GenerateTokenPair(handler.db, user)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(tokens)
}
