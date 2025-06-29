package main

import (
	"database/sql"
	"github.com/SwanHtetAungPhyo/product-api/database"
	"github.com/SwanHtetAungPhyo/product-api/elasticsearch"
	"github.com/SwanHtetAungPhyo/product-api/handlers"
	"github.com/SwanHtetAungPhyo/product-api/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://localhost:9200"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	postUrl := os.Getenv("POSTGRES_URL")
	db := database.Connect(postUrl)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(db)

	es, err := elasticsearch.NewClient(url)
	if err != nil {
		log.Fatal("Failed to connect to Elasticsearch:", err)
	}

	// Initialize services
	productService := services.NewProductService(db, es)
	productHandler := handlers.NewProductHandler(productService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Routes
	api := app.Group("/api/v1")

	// Product routes
	api.Post("/products", productHandler.CreateProduct)
	api.Get("/products", productHandler.GetAllProducts)
	api.Get("/products/search", productHandler.SearchProducts)
	api.Get("/products/:id", productHandler.GetProduct)
	api.Put("/products/:id", productHandler.UpdateProduct)
	api.Delete("/products/:id", productHandler.DeleteProduct)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "Product API is running",
		})
	})

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
