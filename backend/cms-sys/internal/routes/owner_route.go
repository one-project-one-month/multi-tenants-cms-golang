package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/handler"
)

func SetupOwnerRoutes(app *fiber.App, handler handler.OwnerHandle) {
	ownerRoute := app.Group("/owners")
	ownerRoute.Post("/create", handler.Create)
	ownerRoute.Put("/update/:id", handler.Update)
	ownerRoute.Get("/getAllOwners", handler.GetAll)
	ownerRoute.Get("/:id", handler.GetOwnerByID)
}
