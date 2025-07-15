package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/handler"
)

func SetupRoutes(app *fiber.App, handler handler.AuthHandle) {
	auth := app.Group("/cms-doc/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/login/mfa/:userid", handler.LoginWithMFA)
	auth.Post("/register", handler.Register)
	auth.Post("/logout", handler.Logout)
	auth.Post("/refresh", handler.Refresh)
	auth.Post("/me", handler.GetMe)
	auth.Put("/profile/:id", handler.UpdateUserProfile)
	auth.Post("/mfa/setup/:userid", handler.SetupMFA)
	auth.Post("/mfa/verify/:userid", handler.VerifyMFASetup)
}
