package auth

import (
	"trophy/internal/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func generateJwt(username string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSigningKey())
}

func GenerateTokenPair(db *gorm.DB, user database.User) (*TokenResponse, error) {
	accessToken, err := generateJwt(user.Username, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshTokenString := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	dbToken := database.RefreshToken{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	if err := db.Create(&dbToken).Error; err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
	}, nil
}
