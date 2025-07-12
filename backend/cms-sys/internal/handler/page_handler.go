package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/service"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
)

type PageHandle interface {
	GetAll(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type PageHandler struct {
	service   service.PageService
	validator *validator.Validate
}

var _ PageHandle = (*PageHandler)(nil)

func NewPageHandler(service service.PageService) PageHandle {
	return &PageHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (p *PageHandler) GetAll(c *fiber.Ctx) error {
	var req types.PaginateRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.InternalServerErrorResponse(c, "Invalid query parameters", err.Error())
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	pages, pagination, err := p.service.GetAllPages(&req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to get all pages", err.Error())
	}

	return utils.PaginatedSuccessResponse(c, "All pages retrieved successfully", pages, *pagination)
}
func (p *PageHandler) Create(c *fiber.Ctx) error {
	var req types.PageCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request", err.Error())
	}

	pageResponse, err := p.service.Create(&req)

	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create page", err.Error())
	}

	return utils.CreatedResponse(c, "Page created successfully", pageResponse)
}

func (p *PageHandler) GetById(c *fiber.Ctx) error {
	id := c.Params("id")
	page, err := p.service.GetById(id)

	if err != nil {
		return utils.NotFoundResponse(c, "Page Not Found")
	}

	return utils.SuccessResponse(c, "Page found successfully", page)
}
func (p *PageHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req types.PageUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}
	if err := p.validator.Struct(req); err != nil {
		return utils.BadRequestResponse(c, "Validation Failed!", err.Error())
	}

	pageResponse, err := p.service.UpdateById(id, &req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update page", err.Error())
	}
	return utils.SuccessResponse(c, "Page updated successfully", pageResponse)
}
