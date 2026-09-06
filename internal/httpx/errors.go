package httpx

import (
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "internal server error"

	if fiberErr, ok := err.(*fiber.Error); ok {
		code = fiberErr.Code
		message = fiberErr.Message
	}

	return c.Status(code).JSON(fiber.Map{"error": message})
}
