package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/aws"

	"github.com/multi-tenants-cms-golang/cms-sys/internal/service"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
)

type PageRequestHandle interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type PageRequestHandler struct {
	service   service.PageRequestService
	validator *validator.Validate
	s3Service *aws.S3Service
}

var _ PageRequestHandle = (*PageRequestHandler)(nil)

func NewPageRequestHandler(
	service service.PageRequestService,
	s3 *aws.S3Service,
) PageRequestHandle {
	return &PageRequestHandler{
		service:   service,
		validator: validator.New(),
		s3Service: s3,
	}
}

func (h *PageRequestHandler) Create(c *fiber.Ctx) error {
	var req types.CreatePageRequest
	if form, err := c.MultipartForm(); err != nil {
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
		uploadedURL, err := h.s3Service.UploadFile(req.LogoFile)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to upload logo", err.Error())
		}
		logoURL = uploadedURL
	}

	pageRequestResponse, err := h.service.CreatePageRequest(req,logoURL)
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
