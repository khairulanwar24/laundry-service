package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func Logger(c *fiber.Ctx) error {
	c.Locals("requestid", "some-random-id")
	return c.Next()
}
