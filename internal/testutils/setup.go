package testutils

import (
	"context"
	"fmt"
	"testing"
	"trophy/internal/database"
	authRouter "trophy/internal/router/auth"
	healthRouter "trophy/internal/router/health"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	ctx := context.Background()

	postgresC, err := testcontainers.Run(
		ctx, "postgres:16-alpine",
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		}),
		testcontainers.WithExposedPorts("5432/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForListeningPort("5432/tcp"),
			),
		),
	)
	testcontainers.CleanupContainer(t, postgresC)
	require.NoError(t, err)

	host, err := postgresC.Host(ctx)
	require.NoError(t, err)
	port, err := postgresC.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port(),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	database.MigrateDatabases(db)
	return db
}

func SetupApp(t *testing.T) (*fiber.App, *gorm.DB) {
	app := fiber.New()
	db := SetupTestDB(t)

	api := app.Group("/api")
	authRouter.Register(api, db)
	healthRouter.Register(api, db)

	return app, db
}
