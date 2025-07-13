package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/handler"
)

func SetupPageRoutes(app *fiber.App, handler handler.PageHandle) {
	route := app.Group("cms/pages")
	route.Get("/", handler.GetAll)
	route.Post("/", handler.Create)
	route.Get("/:id", handler.GetById)
	route.Put("/:id", handler.Update)
}
