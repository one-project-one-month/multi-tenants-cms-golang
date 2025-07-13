package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/handler"
)

func SetupOwnerRoutes(app *fiber.App, handler handler.OwnerHandle) {
	ownerRoute := app.Group("cms/owners")
	ownerRoute.Post("/", handler.Create)
	ownerRoute.Put("/:id", handler.Update)
	ownerRoute.Get("/", handler.GetAll)
	ownerRoute.Get("/:id", handler.GetOwnerByID)
	ownerRoute.Delete("/", handler.Delete)
}
