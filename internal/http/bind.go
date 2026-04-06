package http

import (
	"bytes"
	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

func Bind[T any](ctx fiber.Ctx, body *T) error {
	decoder := json.NewDecoder(bytes.NewReader(ctx.Body()))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}
	return nil
}
