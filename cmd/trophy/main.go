package main

import (
	"time"
	"trophy/internal/database"
	authRouter "trophy/internal/router/auth"
	clipRouter "trophy/internal/router/clips"
	healthRouter "trophy/internal/router/health"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
)

// @title						Trophy API
// @version					1.0
// @description				Self-hostable backend API for Trophy.
// @host						localhost:3000
// @BasePath					/api
// @schemes					http
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found!")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := database.MigrateDatabases(db); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:        40,
		Expiration: 30 * time.Second,
	}))

	app.Get("/docs/*", swaggo.HandlerDefault)

	api := app.Group("/api")

	authRouter.Register(api, db)
	healthRouter.Register(api, db)
	clipRouter.Register(api, db)

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
