package routes

import (
	"github.com/gofiber/fiber/v2"
)

func DummyRoutes(app *fiber.App) {
	app.Get("/dummy", func(c *fiber.Ctx) error {
		return c.SendString("Dummy server")
	})
}
