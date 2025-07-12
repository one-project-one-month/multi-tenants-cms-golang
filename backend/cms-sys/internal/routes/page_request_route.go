package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/handler"
)

func SetupPageRequestRoutes(app *fiber.App, handler handler.PageRequestHandle) {
	pageRequest := app.Group("/page-request")
	pageRequest.Post("/", handler.Create)
	pageRequest.Get("/", handler.GetAll)
}
