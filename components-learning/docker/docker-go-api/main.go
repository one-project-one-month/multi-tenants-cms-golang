package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	go func() {
		log.Printf("Listening on port %s", port)
		err := app.Listen(":" + port)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}()

	osCh := make(chan os.Signal, 1)
	signal.Notify(osCh, os.Interrupt)
	<-osCh
	log.Println("Shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Fatal(err.Error())
	}
}
