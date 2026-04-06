package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;size:255"`
	Password string `gorm:"not null;size:255"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"unique;not null"`
	UserID    uint      `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type Clip struct {
	gorm.Model
	Title     string `gorm:"not null"`
	VideoHash string `gorm:"unique;not null"`
	UserID    uint   `gorm:"not null"`
}

type Comment struct {
	gorm.Model
	Message string `gorm:"not null"`
	Likes   int    `gorm:"default:0"`
	ClipID  uint   `gorm:"not null"`
	UserID  uint   `gorm:"not null"`
}

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}

		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}

		user := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")

		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host,
			port,
			user,
			password,
			name,
		)
	}

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func MigrateDatabases(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &RefreshToken{}, &Clip{}, &Comment{})
}
