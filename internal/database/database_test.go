package database_test

import (
	"testing"
	"trophy/internal/database"
	"trophy/internal/dbtest"

	"github.com/stretchr/testify/require"
)

func TestMigrateDatabases(t *testing.T) {
	db := dbtest.SetupDB(t)

	err := database.MigrateDatabases(db)
	require.NoError(t, err)
}

func TestModelsCanBePersisted(t *testing.T) {
	db := dbtest.SetupDB(t)

	user := database.User{Username: "db-test-user", Password: "password"}
	require.NoError(t, db.Create(&user).Error)

	clip := database.Clip{Title: "test clip", VideoHash: "abcdef1234567890", UserID: user.ID}
	require.NoError(t, db.Create(&clip).Error)

	comment := database.Comment{Message: "great clip", ClipID: clip.ID, UserID: user.ID}
	require.NoError(t, db.Create(&comment).Error)

	refreshToken := database.RefreshToken{Token: "refresh-token", UserID: user.ID}
	require.NoError(t, db.Create(&refreshToken).Error)
}
