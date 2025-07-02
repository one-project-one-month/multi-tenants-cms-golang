package main

import (
	"github.com/multi-tenants-cms-golang/cms-sys/internal/handler"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/repository"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type DISection struct {
	repo    repository.AuthRepository
	srv     service.AuthService
	handler handler.AuthHandle
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Go Fiber App",
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello Fiber!",
			"status":  "success",
		})
	})
	//
	//routes.SetupRoutes(app,)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	log.Println("🚀 Server starting on :8080")
	log.Fatal(app.Listen(":8080"))
}

//func DependencyInjectionSection()
