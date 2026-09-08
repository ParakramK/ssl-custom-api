package router

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
)

func apiKeyAuth(expectedKey string) fiber.Handler {
	return keyauth.New(keyauth.Config{
		Validator: func(_ fiber.Ctx, key string) (bool, error) {
			if expectedKey == "" || key == "" {
				return false, nil
			}
			return subtle.ConstantTimeCompare([]byte(key), []byte(expectedKey)) == 1, nil
		},
		Extractor: extractors.Chain(
			extractors.FromHeader("X-API-Key"),
		),
		ErrorHandler: func(c fiber.Ctx, _ error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or invalid API key",
			})
		},
	})
}
