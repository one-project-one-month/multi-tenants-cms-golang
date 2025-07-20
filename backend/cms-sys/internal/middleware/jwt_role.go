package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
)

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user_claims").(*utils.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No user claims found",
			})
		}

		for _, role := range roles {
			if claims.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}

func RequireTokenType(tokenType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user_claims").(*utils.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No user claims found",
			})
		}

		if claims.TokenType != tokenType {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token type",
			})
		}

		return c.Next()
	}
}

func RequireAccessToken() fiber.Handler {
	return RequireTokenType("access")
}

func RequireOwnership() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User ID not found in token",
			})
		}

		return c.Next()
	}
}
