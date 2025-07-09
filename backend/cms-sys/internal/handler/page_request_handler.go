package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

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
}

var _ PageRequestHandle = (*PageRequestHandler)(nil)

func NewPageRequestHandler(service service.PageRequestService) PageRequestHandle {
	return &PageRequestHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *PageRequestHandler) Create(c *fiber.Ctx) error {
	var req types.CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.BadRequestResponse(c, "Validation failed", err.Error())
	}

	pageRequestResponse, err := h.service.CreatePageRequest(req)
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
