package health

import (
	"trophy/internal/database"

	"github.com/gofiber/fiber/v3"
)

type healthChecks struct {
	Database string `json:"database"`
}

type healthResponse struct {
	Status string       `json:"status"`
	Checks healthChecks `json:"checks"`
}

// Health returns service health and dependency checks.
// @Summary Health check
// @Description Check if the API and database are healthy
// @Tags Health
// @Success 200 {object} healthResponse
// @Failure 503 {object} healthResponse
// @Router /health [get]
func Health(ctx fiber.Ctx) error {
	db, err := database.DB.DB()
	if err != nil || db.Ping() != nil {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "unhealthy",
			"checks": fiber.Map{
				"database": "unreachable",
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "healthy",
		"checks": fiber.Map{
			"database": "ok",
		},
	})
}
