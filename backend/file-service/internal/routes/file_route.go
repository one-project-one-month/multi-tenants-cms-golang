package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/file-service/internal/handler"
)

func FileRoute(app *fiber.App, handler handler.FileSystemHandler) {
	api := app.Group("/api/v1")

	api.Post("/files", handler.UploadHandler)
	api.Get("/files/*", handler.GetHandler)
	api.Put("/files/*", handler.UpdateHandler)
	api.Delete("/files/*", handler.DeleteHandler)
	api.Get("/files", handler.ListHandler)
}
