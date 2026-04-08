package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestBind(t *testing.T) {
	type requestBody struct {
		Title string `json:"title"`
	}

	app := fiber.New()
	app.Post("/bind", func(c fiber.Ctx) error {
		var body requestBody
		if err := Bind(c, &body); err != nil {
			return nil
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "valid json body",
			body:           `{"title":"clip"}`,
			expectedStatus: fiber.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/bind", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
