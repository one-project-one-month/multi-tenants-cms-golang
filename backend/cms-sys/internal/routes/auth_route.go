package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/handler"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/middleware"
	"os"
)

func SetupRoutes(app *fiber.App, handler handler.AuthHandle) {
	auth := app.Group("/cms/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Post("/verify-email", handler.VerifyEmail)
	auth.Post("/mfa/setup/:userid", handler.SetupMFA)
	auth.Post("/mfa/verify/:userid", handler.VerifyMFASetup)

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	auth.Post("/refresh",
		middleware.JWTMiddleware(jwtSecret),
		middleware.RequireTokenType("refresh"),
		handler.Refresh,
	)

	protected := auth.Group("",
		middleware.JWTMiddleware(jwtSecret),
		middleware.RequireAccessToken(),
	)

	protected.Post("/logout", handler.Logout)
	protected.Get("/me", handler.GetMe)
	protected.Put("/profile/:id", handler.UpdateUserProfile)
	protected.Post("/mfa/login/:userid", handler.LoginWithMFA)
}
