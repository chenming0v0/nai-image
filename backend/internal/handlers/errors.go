package handlers

import (
	"errors"

	"naiimage/backend/internal/nai"

	"github.com/gofiber/fiber/v2"
)

func fiberError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

func upstreamFiberError(c *fiber.Ctx, err error) error {
	var ue *nai.UpstreamError
	if errors.As(err, &ue) {
		status := ue.Status
		if status < 400 || status > 599 {
			status = 502
		}
		return fiberError(c, status, ue.Message)
	}
	return fiberError(c, 502, err.Error())
}
