package http

import "github.com/gofiber/fiber/v3"

func Bind[T any](ctx fiber.Ctx, body *T) error {
	if err := ctx.Bind().Body(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}
	return nil
}
