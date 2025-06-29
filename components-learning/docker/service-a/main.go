package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"log"
	"mime/multipart"
	"os"
	"os/signal"
)

type FileUploadMessage struct {
	FileName    string `json:"file_name"`
	FileContent []byte `json:"file_content"`
	BucketName  string `json:"bucket_name"`
	ObjectKey   string `json:"object_key"`
}

var natsConn *nats.Conn

func main() {
	var err error
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://nats:4222"
	}

	natsConn, err = nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer natsConn.Close()
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
	app.Post("/upload", uploadHandler)

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

func uploadHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "no file attachment found",
		})
	}

	fileContent, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "no file attachment found",
		})
	}
	defer func(fileContent multipart.File) {
		err := fileContent.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err.Error())
		}
	}(fileContent)

	content := make([]byte, file.Size)
	_, err = fileContent.Read(content)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "cannot read file",
		})
	}

	message := &FileUploadMessage{
		FileName:    file.Filename,
		FileContent: content,
		BucketName:  os.Getenv("BUCKET_NAME"),
		ObjectKey:   file.Filename,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "cannot marshal message",
		})
	}
	err = natsConn.Publish("file.upload", messageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "cannot publish message to nats server",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "file uploaded request send",
	})
}
