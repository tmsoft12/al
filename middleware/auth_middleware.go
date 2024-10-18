// middleware/auth_middleware.go
package middleware

import (
	"rr/service"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("jwt")
		if cookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Şahsyýet ýok"})
		}

		token, err := service.ValidateToken(cookie)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Şahsyýet ýok"})
		}

		return c.Next()
	}
}
