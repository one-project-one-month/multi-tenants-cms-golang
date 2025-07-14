package handler

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/file-service/internal/service"
	"github.com/multi-tenants-cms-golang/file-service/internal/types"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"path/filepath"
	"time"
)

type (
	FileSystemHandler interface {
		UploadHandler(ctx *fiber.Ctx) error
		GetHandler(ctx *fiber.Ctx) error
		ListHandler(ctx *fiber.Ctx) error
		UpdateHandler(ctx *fiber.Ctx) error
		DeleteHandler(ctx *fiber.Ctx) error
	}

	FileSystemHandlerImpl struct {
		fh  *logrus.Logger
		srv service.FileSystemService
	}
)

var _ FileSystemHandler = (*FileSystemHandlerImpl)(nil)

func NewFileSystemHandler(
	log *logrus.Logger,
	srv service.FileSystemService,
) *FileSystemHandlerImpl {
	return &FileSystemHandlerImpl{
		fh:  log,
		srv: srv,
	}
}

func (f FileSystemHandlerImpl) UploadHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
	}

	fileContent, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read file"})
	}
	defer func(fileContent multipart.File) {
		err := fileContent.Close()
		if err != nil {
			f.fh.Error(err.Error())
		}
	}(fileContent)

	data := make([]byte, file.Size)
	_, err = fileContent.Read(data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read file content"})
	}

	options := &types.UploadOptions{
		Filename:    file.Filename,
		Directory:   c.FormValue("directory", "uploads"),
		ContentType: file.Header.Get("Content-Type"),
		Overwrite:   c.FormValue("overwrite") == "true",
		Metadata: map[string]string{
			"uploaded_by": c.FormValue("user_id", "anonymous"),
			"upload_time": time.Now().Format(time.RFC3339),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fileInfo, err := f.srv.UploadFile(ctx, data, options)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "File uploaded successfully",
		"file":    fileInfo,
	})
}

func (f FileSystemHandlerImpl) GetHandler(c *fiber.Ctx) error {
	path := c.Params("*")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	data, err := f.srv.DownloadFile(ctx, path)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}

	filename := filepath.Base(path)
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return c.Send(data)
}

func (f FileSystemHandlerImpl) ListHandler(c *fiber.Ctx) error {
	directory := c.Query("directory", ".")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	files, err := f.srv.ListFiles(ctx, directory)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"directory": directory,
		"files":     files,
	})
}

func (f FileSystemHandlerImpl) UpdateHandler(c *fiber.Ctx) error {
	path := c.Params("*")

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
	}

	fileContent, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read file"})
	}
	defer func(fileContent multipart.File) {
		err := fileContent.Close()
		if err != nil {
			f.fh.Error(err.Error())
		}
	}(fileContent)

	data := make([]byte, file.Size)
	_, err = fileContent.Read(data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read file content"})
	}

	options := &types.UploadOptions{
		Filename:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Metadata: map[string]string{
			"updated_by":  c.FormValue("user_id", "anonymous"),
			"update_time": time.Now().Format(time.RFC3339),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fileInfo, err := f.srv.UpdateFile(ctx, path, data, options)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "File updated successfully",
		"file":    fileInfo,
	})
}

func (f FileSystemHandlerImpl) DeleteHandler(c *fiber.Ctx) error {
	path := c.Params("*")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := f.srv.DeleteFile(ctx, path)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "File deleted successfully"})
}
