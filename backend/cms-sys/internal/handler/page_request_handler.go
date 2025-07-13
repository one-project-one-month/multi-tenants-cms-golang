package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/aws"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/file_srv_communication"
	"io"
	"mime/multipart"
	"os"

	"github.com/multi-tenants-cms-golang/cms-sys/internal/service"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
)

type PageRequestHandle interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type PageRequestHandler struct {
	service    service.PageRequestService
	validator  *validator.Validate
	s3Service  *aws.S3Service
	fileClient *file_srv_communication.FileClient
}

var _ PageRequestHandle = (*PageRequestHandler)(nil)

func NewPageRequestHandler(
	service service.PageRequestService,
	s3 *aws.S3Service,
	fileClient *file_srv_communication.FileClient,
) PageRequestHandle {
	return &PageRequestHandler{
		service:    service,
		validator:  validator.New(),
		s3Service:  s3,
		fileClient: fileClient,
	}
}

func (h *PageRequestHandler) Create(c *fiber.Ctx) error {
	var req types.CreatePageRequest

	if form, err := c.MultipartForm(); err != nil {
		if form == nil {
			return utils.BadRequestResponse(c, "Invalid request body", err.Error())
		}
		if files := form.File["logo"]; len(files) > 0 {
			req.LogoFile = files[0]
		}
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.BadRequestResponse(c, "Validation failed", err.Error())
	}

	var logoURL *string
	if req.LogoFile != nil {
		tempFile, err := os.CreateTemp("", "logo-*.tmp")
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to create temp file", err.Error())
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				fmt.Printf("failed to remove file: %v\n", err)
			}
		}(tempFile.Name())
		defer func(tempFile *os.File) {
			err := tempFile.Close()
			if err != nil {
				fmt.Printf("failed to close file: %v\n", err)
			}
		}(tempFile)

		fileContent, err := req.LogoFile.Open()
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to open uploaded file", err.Error())
		}
		defer func(fileContent multipart.File) {
			err := fileContent.Close()
			if err != nil {
				fmt.Printf("failed to close file: %v\n", err)
			}
		}(fileContent)

		if _, err := io.Copy(tempFile, fileContent); err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to save uploaded file", err.Error())
		}

		uploadReq := types.FileUploadRequest{
			FilePath:   tempFile.Name(),
			Directory:  "page-requests/logos",
			UserID:     req.OwnerID,
			Overwrite:  false,
			ServiceURL: "http://localhost:9009/file",
			Metadata: map[string]string{
				"purpose":      "page_request_logo",
				"page_request": *req.PageUrl,
			},
		}

		uploadResp, err := h.fileClient.UploadFile(c.Context(), uploadReq)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to upload logo to file service", err.Error())
		}

		if fileInfo, ok := uploadResp.File.(map[string]interface{}); ok {
			if url, ok := fileInfo["url"].(string); ok {
				logoURL = &url
			}
		}
	}

	pageRequestResponse, err := h.service.CreatePageRequest(req, logoURL)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create page request", err.Error())
	}

	return utils.CreatedResponse(c, "Page request created", pageRequestResponse)
}

func (h *PageRequestHandler) GetAll(c *fiber.Ctx) error {
	var req types.PaginateRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	pageRequests, pagination, err := h.service.GetAllPageRequests(&req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to get page requests", err.Error())
	}

	return utils.PaginatedSuccessResponse(c, "Page requests retrieved successfully", pageRequests, *pagination)
}
