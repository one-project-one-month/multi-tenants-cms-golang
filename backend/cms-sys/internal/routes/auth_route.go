package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/handler"
)

func SetupRoutes(
	app *fiber.App,
	handler handler.AuthHandle,
) {
	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Post("/logout", handler.Logout)
	auth.Post("/refresh", handler.Refresh)
}
