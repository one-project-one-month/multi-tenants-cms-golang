package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
)

func GetUserClaims(c *fiber.Ctx) (*utils.Claims, bool) {
	claims, ok := c.Locals("user_claims").(*utils.Claims)
	return claims, ok
}
