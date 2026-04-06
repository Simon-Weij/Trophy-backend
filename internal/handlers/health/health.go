package health

import (
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
func (handler *Handler) Health(ctx fiber.Ctx) error {
	sqlDB, err := handler.db.DB()
	if err != nil || sqlDB.Ping() != nil {
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

// Protected returns 200 OK for testing protected endpoints.
// @Summary Protected test endpoint
// @Description Test endpoint that requires authentication
// @Tags Health
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /health/protected [get]
func (handler *Handler) Protected(ctx fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
