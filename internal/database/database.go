package database

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v3/log"
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

var DB *gorm.DB

func MigrateDatabases() {
	var err error

	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("SSL_MODE")

	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=5432 sslmode=%s", user, password, dbname, sslmode)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Couldn't connect to database %v", err)
	}

	DB.AutoMigrate(&User{}, &RefreshToken{}, &Clip{}, &Comment{})
}
