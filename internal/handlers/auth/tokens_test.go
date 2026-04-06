package auth

import (
	"testing"
	"time"
	"trophy/internal/database"
	"trophy/internal/dbtest"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func parseTokenWithKey(tokenString string, keyFunc jwt.Keyfunc, opts ...jwt.ParserOption) (*jwt.Token, error) {
	return jwt.NewParser(opts...).ParseWithClaims(tokenString, &Claims{}, keyFunc)
}

func parseTokenWithValidation(tokenString string) (*jwt.Token, error) {
	return parseTokenWithKey(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSigningKey(), nil
	})
}

func parseTokenWithoutValidation(tokenString string) (*jwt.Token, error) {
	return parseTokenWithKey(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSigningKey(), nil
	}, jwt.WithoutClaimsValidation())
}

func Test_generateJwt(t *testing.T) {
	t.Run("generates valid token", func(t *testing.T) {
		username := "testuser"
		token, err := generateJwt(username, time.Second*5)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("token expires after specified duration", func(t *testing.T) {
		jwtUsername := "username"
		tokenValidTime := time.Second * 1
		token, err := generateJwt(jwtUsername, tokenValidTime)
		require.NoError(t, err)

		parsedToken, err := parseTokenWithValidation(token)
		require.NoError(t, err)
		require.NotNil(t, parsedToken)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(*Claims)
		require.True(t, ok)
		assert.Equal(t, jwtUsername, claims.Username)

		time.Sleep(tokenValidTime + time.Millisecond*101)

		expiredToken, err := parseTokenWithoutValidation(token)
		require.NoError(t, err)
		require.NotNil(t, expiredToken)

		expiredClaims, ok := expiredToken.Claims.(*Claims)
		require.True(t, ok)

		expirationTime, err := expiredClaims.GetExpirationTime()
		require.NoError(t, err)
		assert.True(t, expirationTime.Before(time.Now()))
	})

	t.Run("username is encoded in claims", func(t *testing.T) {
		jwtUsername := "testuser"
		token, err := generateJwt(jwtUsername, time.Second*60)
		require.NoError(t, err)

		parsedToken, err := parseTokenWithValidation(token)
		require.NoError(t, err)

		claims, ok := parsedToken.Claims.(*Claims)
		require.True(t, ok)
		assert.Equal(t, jwtUsername, claims.Username)
	})

	t.Run("token is invalid with wrong signing key", func(t *testing.T) {
		jwtUsername := "username"
		token, err := generateJwt(jwtUsername, time.Second*60)
		require.NoError(t, err)

		wrongKey := []byte("wrong-secret-key")
		parsedToken, err := jwt.NewParser().ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return wrongKey, nil
		})

		assert.Error(t, err)
		assert.False(t, parsedToken.Valid)
	})

	t.Run("expiration time is approximately duration away", func(t *testing.T) {
		jwtUsername := "username"
		duration := time.Second * 30
		beforeGeneration := time.Now()
		token, err := generateJwt(jwtUsername, duration)
		afterGeneration := time.Now()
		require.NoError(t, err)

		parsedToken, err := parseTokenWithValidation(token)
		require.NoError(t, err)

		claims, ok := parsedToken.Claims.(*Claims)
		require.True(t, ok)

		expirationTime, err := claims.GetExpirationTime()
		require.NoError(t, err)

		expectedMinTime := beforeGeneration.Add(duration)
		expectedMaxTime := afterGeneration.Add(duration + 2*time.Second)

		assert.True(t,
			expirationTime.Time.After(expectedMinTime.Add(-time.Second)) &&
				expirationTime.Time.Before(expectedMaxTime),
			"Expiration time should be approximately %v away from generation time",
			duration)
	})
}

func TestGenerateTokenPair(t *testing.T) {
	t.Run("generates valid token pair", func(t *testing.T) {
		db := dbtest.SetupDB(t)

		user := database.User{Username: "testuser", Password: "hashedpwd"}
		result := db.Create(&user)
		require.NoError(t, result.Error)

		tokenResponse, err := GenerateTokenPair(db, user)
		require.NoError(t, err)
		require.NotNil(t, tokenResponse)

		assert.NotEmpty(t, tokenResponse.AccessToken)
		assert.NotEmpty(t, tokenResponse.RefreshToken)
	})

	t.Run("access token is valid JWT with user claims", func(t *testing.T) {
		db := dbtest.SetupDB(t)

		user := database.User{Username: "jwtuser", Password: "hashedpwd"}
		result := db.Create(&user)
		require.NoError(t, result.Error)

		tokenResponse, err := GenerateTokenPair(db, user)
		require.NoError(t, err)

		parsedToken, err := parseTokenWithValidation(tokenResponse.AccessToken)
		require.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(*Claims)
		require.True(t, ok)
		assert.Equal(t, user.Username, claims.Username)
	})

	t.Run("access token has 15 minute expiration", func(t *testing.T) {
		db := dbtest.SetupDB(t)

		user := database.User{Username: "expiryuser", Password: "hashedpwd"}
		result := db.Create(&user)
		require.NoError(t, result.Error)

		beforeGeneration := time.Now()
		tokenResponse, err := GenerateTokenPair(db, user)
		afterGeneration := time.Now()
		require.NoError(t, err)

		parsedToken, err := parseTokenWithValidation(tokenResponse.AccessToken)
		require.NoError(t, err)

		claims, ok := parsedToken.Claims.(*Claims)
		require.True(t, ok)

		expirationTime, err := claims.GetExpirationTime()
		require.NoError(t, err)

		expectedMinTime := beforeGeneration.Add(15 * time.Minute)
		expectedMaxTime := afterGeneration.Add(15*time.Minute + 2*time.Second)

		assert.True(t,
			expirationTime.Time.After(expectedMinTime.Add(-time.Second)) &&
				expirationTime.Time.Before(expectedMaxTime),
			"Access token should expire in approximately 15 minutes")
	})

	t.Run("refresh token is stored in database", func(t *testing.T) {
		db := dbtest.SetupDB(t)

		user := database.User{Username: "dbuser", Password: "hashedpwd"}
		result := db.Create(&user)
		require.NoError(t, result.Error)

		tokenResponse, err := GenerateTokenPair(db, user)
		require.NoError(t, err)

		var storedToken database.RefreshToken
		result = db.Where("token = ?", tokenResponse.RefreshToken).First(&storedToken)
		require.NoError(t, result.Error)

		assert.Equal(t, user.ID, storedToken.UserID)
		assert.Equal(t, tokenResponse.RefreshToken, storedToken.Token)
	})

	t.Run("refresh token has 7 day expiration", func(t *testing.T) {
		db := dbtest.SetupDB(t)

		user := database.User{Username: "refreshuser", Password: "hashedpwd"}
		result := db.Create(&user)
		require.NoError(t, result.Error)

		beforeGeneration := time.Now()
		tokenResponse, err := GenerateTokenPair(db, user)
		afterGeneration := time.Now()
		require.NoError(t, err)

		var storedToken database.RefreshToken
		result = db.Where("token = ?", tokenResponse.RefreshToken).First(&storedToken)
		require.NoError(t, result.Error)

		expectedMinTime := beforeGeneration.Add(7 * 24 * time.Hour)
		expectedMaxTime := afterGeneration.Add(7*24*time.Hour + 2*time.Second)

		assert.True(t,
			storedToken.ExpiresAt.After(expectedMinTime.Add(-time.Second)) &&
				storedToken.ExpiresAt.Before(expectedMaxTime),
			"Refresh token should expire in approximately 7 days")
	})

	t.Run("refresh token is unique UUID", func(t *testing.T) {
		db := dbtest.SetupDB(t)

		user := database.User{Username: "uuiduser", Password: "hashedpwd"}
		result := db.Create(&user)
		require.NoError(t, result.Error)

		tokenResponse1, err := GenerateTokenPair(db, user)
		require.NoError(t, err)

		tokenResponse2, err := GenerateTokenPair(db, user)
		require.NoError(t, err)

		assert.NotEqual(t, tokenResponse1.RefreshToken, tokenResponse2.RefreshToken)
	})

	t.Run("returns error when user creation in database fails", func(t *testing.T) {
		var nilDB *gorm.DB

		user := database.User{Username: "erroruser", Password: "hashedpwd"}
		tokenResponse, err := GenerateTokenPair(nilDB, user)

		assert.Error(t, err)
		assert.Nil(t, tokenResponse)
	})
}
